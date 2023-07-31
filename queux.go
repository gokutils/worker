package worker

import (
	"context"
	"time"

	"github.com/gokutils/container/queux"
)

type WorkerQeux[T any] interface {
	Do(ctx context.Context, task T) error
}

type params struct {
	Timeout time.Duration
	OnError func(err error)
}

type Queux[T any] struct {
	coworker      int
	worker        WorkerQeux[T]
	golbalCtx     context.Context
	managerCtx    context.Context
	managerCancel func()
	queux         *queux.QueuxMemory[T]
	params        params
}

func (impl *Queux[T]) do(msg T) {
	var ctx context.Context
	var cancel func()
	if impl.params.Timeout > 0 {
		ctx, cancel = context.WithTimeout(impl.managerCtx, impl.params.Timeout)
	} else {
		ctx, cancel = context.WithCancel(impl.managerCtx)
	}
	defer cancel()
	if err := impl.worker.Do(ctx, msg); err != nil {
		if impl.params.OnError != nil {
			impl.params.OnError(err)
		}
	}
}

func (impl *Queux[T]) start() {
	go func() {
		for {
			v, err := impl.queux.GetOrWait(impl.managerCtx)
			if err != nil {
				return
			}
			impl.do(v)
		}
	}()
}

func (impl *Queux[T]) Start() {
	if impl.managerCancel != nil {
		return
	}
	impl.managerCtx, impl.managerCancel = context.WithCancel(impl.golbalCtx)
	for i := 0; i < impl.coworker; i++ {
		impl.start()
	}
}

func (impl *Queux[T]) Push(v T) error {
	return impl.queux.Push(v)
}

func (impl *Queux[T]) Stop() {
	if impl.managerCancel == nil {
		return
	}
	impl.managerCancel()
	impl.managerCancel = nil
	impl.managerCtx = nil
}

func NewQueux[T any](ctx context.Context, num int, Worker WorkerQeux[T], opts ...Option) *Queux[T] {
	params := params{}
	for i := range opts {
		opts[i](&params)
	}

	return &Queux[T]{
		coworker:  num,
		golbalCtx: ctx,
		worker:    Worker,
		queux:     queux.NewQueuxMemory[T](),
		params:    params,
	}
}
