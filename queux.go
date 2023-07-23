package worker

import (
	"context"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
)

type WorkerQeux[T any] interface {
	Do(ctx context.Context, task T) error
}

type params struct {
	Timeout time.Duration
	OnError func(err error)
}

type Qeux[T any] struct {
	coworker      int
	worker        WorkerQeux[T]
	golbalCtx     context.Context
	managerCtx    context.Context
	managerCancel func()
	queux         *goconcurrentqueue.FIFO
	params        params
}

func (impl *Qeux[T]) start() {
	go func() {
		for {
			v, err := impl.queux.DequeueOrWaitForNextElementContext(impl.managerCtx)
			if err != nil {
				return
			}
			var ctx context.Context
			var cancel func()
			if impl.params.Timeout > 0 {
				ctx, cancel = context.WithTimeout(impl.managerCtx, impl.params.Timeout)
			} else {
				ctx, cancel = context.WithCancel(impl.managerCtx)
			}
			defer cancel()
			if err := impl.worker.Do(ctx, v.(T)); err != nil {
				if impl.params.OnError != nil {
					impl.params.OnError(err)
				}
			}
		}
	}()
}

func (impl *Qeux[T]) Start() {
	if impl.managerCancel != nil {
		return
	}
	impl.managerCtx, impl.managerCancel = context.WithCancel(impl.golbalCtx)
	for i := 0; i < impl.coworker; i++ {
		impl.start()
	}
}
func (impl *Qeux[T]) Push(v T) error {
	return impl.queux.Enqueue(v)
}

func (impl *Qeux[T]) Stop() {
	if impl.managerCancel == nil {
		return
	}
	impl.managerCancel()
	impl.managerCancel = nil
	impl.managerCtx = nil
}

func NewQueux[T any](ctx context.Context, num int, Worker WorkerQeux[T], opts ...Option) *Qeux[T] {
	params := params{}
	for i := range opts {
		opts[i](&params)
	}

	return &Qeux[T]{
		coworker:  num,
		golbalCtx: ctx,
		worker:    Worker,
		queux:     goconcurrentqueue.NewFIFO(),
		params:    params,
	}
}
