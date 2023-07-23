package worker

import (
	"context"
	"time"

	"log"

	"github.com/enriquebris/goconcurrentqueue"
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

	queux     *goconcurrentqueue.FIFO
	askedTick chan time.Time

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
	ticker := time.NewTicker(impl.tickTime)
	go func() {
		for {
			select {
			case _, ok := <-ticker.C:
				if !ok {
					ticker.Stop()
					return
				}
				impl.do()
			case _, ok := <-impl.askedTick:
				if !ok {
					ticker.Stop()
					return
				}
				impl.do()
			case <-impl.managerCtx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (impl *Tick) Start() {
	if impl.managerCancel != nil {
		return
	}
	impl.managerCtx, impl.managerCancel = context.WithCancel(impl.golbalCtx)
	impl.askedTick = make(chan time.Time)
	for i := 0; i < impl.coworker; i++ {
		impl.start()
	}
	go func() {
		for {
			_, err := impl.queux.DequeueOrWaitForNextElementContext(impl.managerCtx)
			if err != nil {
				log.Printf("Error on dequeux tick %s", err.Error())
				impl.Stop()
				return
			}
			impl.askedTick <- time.Now()
		}
	}()
}

func (impl *Tick) TickNow() error {
	return impl.queux.Enqueue(time.Now())
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
		queux:     goconcurrentqueue.NewFIFO(),
		params:    params,
	}
}
