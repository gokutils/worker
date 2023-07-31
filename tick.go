package worker

import (
	"context"
	"time"
)

type WorkerTick interface {
	Do(context.Context) error
}

type Tick struct {
	golbalCtx     context.Context
	managerCtx    context.Context
	managerCancel func()

	coworker int
	worker   WorkerTick

	askedTick chan struct{}

	tickTime time.Duration

	params params
}

func (impl *Tick) do() {
	var ctx context.Context
	var cancel func()
	if impl.params.Timeout > 0 {
		ctx, cancel = context.WithTimeout(impl.managerCtx, impl.params.Timeout)
	} else {
		ctx, cancel = context.WithCancel(impl.managerCtx)
	}
	defer cancel()
	if err := impl.worker.Do(ctx); err != nil {
		if impl.params.OnError != nil {
			impl.params.OnError(err)
		}
	}
}

func (impl *Tick) start() {
	go func() {
		impl.do()
		for {
			if impl.managerCtx == nil {
				return
			}
			select {
			case _, ok := <-time.After(impl.tickTime):
				if !ok {
					return
				}
				impl.do()
			case _, ok := <-impl.askedTick:
				if !ok {
					return
				}
				impl.do()
			case <-impl.managerCtx.Done():
				return
			}
		}
	}()

}

func (impl *Tick) TickNow() error {
	if impl.managerCancel == nil {
		return nil
	}
	select {
	case impl.askedTick <- struct{}{}:
		return nil
	default:
		return nil
	}
}

func (impl *Tick) Start() {
	if impl.managerCancel != nil {
		return
	}
	impl.managerCtx, impl.managerCancel = context.WithCancel(impl.golbalCtx)
	impl.askedTick = make(chan struct{})
	for i := 0; i < impl.coworker; i++ {
		impl.start()
	}
}

func (impl *Tick) Stop() {
	if impl.managerCancel == nil {
		return
	}
	impl.managerCancel()
	impl.managerCancel = nil
	impl.managerCtx = nil
}

func NewTick(ctx context.Context, num int, tickTime time.Duration, Worker WorkerTick, opts ...Option) *Tick {
	params := params{}
	for i := range opts {
		opts[i](&params)
	}

	return &Tick{
		coworker:  num,
		golbalCtx: ctx,
		tickTime:  tickTime,
		worker:    Worker,
		params:    params,
	}
}
