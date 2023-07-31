package worker

import (
	"context"
	"time"

	"github.com/gokutils/container/queux"
)

type backgroundCtxKey string

const (
	backgroundCtxKeyStart backgroundCtxKey = "background-task-start"
)

type Task[T any, R any] interface {
	Do(ctx context.Context, task T) (R, error)
}

type Responce[R any] struct {
	Responce R
	Error    error
}

type MessageWithResponce[T any, R any] struct {
	Ctx    context.Context
	Input  T
	Output chan Responce[R]
}

type Background[T any, R any] struct {
	coworker      int
	worker        Task[T, R]
	golbalCtx     context.Context
	managerCtx    context.Context
	managerCancel func()
	queux         *queux.QueuxMemory[*MessageWithResponce[T, R]]
	params        params
}

func (impl *Background[T, R]) do(msg *MessageWithResponce[T, R]) {
	var ctx context.Context
	var cancel func()
	if impl.params.Timeout > 0 {
		ctx, cancel = context.WithTimeout(msg.Ctx, impl.params.Timeout)
	} else {
		ctx, cancel = context.WithCancel(msg.Ctx)
	}
	defer cancel()
	resp, err := impl.worker.Do(ctx, msg.Input)
	if err != nil {
		if impl.params.OnError != nil {
			impl.params.OnError(err)
		}

	}
	msg.Output <- Responce[R]{
		Responce: resp,
		Error:    err,
	}
}

func (impl *Background[T, R]) start() {
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

func (impl *Background[T, R]) Start() {
	if impl.managerCancel != nil {
		return
	}
	impl.managerCtx, impl.managerCancel = context.WithCancel(impl.golbalCtx)
	for i := 0; i < impl.coworker; i++ {
		impl.start()
	}
}

func (impl *Background[T, R]) Do(ctx context.Context, v T) (R, error) {
	tmp := &MessageWithResponce[T, R]{
		Ctx:    context.WithValue(ctx, backgroundCtxKeyStart, time.Now()),
		Input:  v,
		Output: make(chan Responce[R]),
	}
	defer close(tmp.Output)
	if err := impl.queux.Push(tmp); err != nil {
		var noop R
		return noop, err
	}
	resp := <-tmp.Output
	return resp.Responce, resp.Error
}

func (impl *Background[T, R]) Stop() {
	if impl.managerCancel == nil {
		return
	}
	impl.managerCancel()
	impl.managerCancel = nil
	impl.managerCtx = nil
}

func NewBackground[T any, R any](ctx context.Context, num int, task Task[T, R], opts ...Option) *Background[T, R] {
	params := params{}
	for i := range opts {
		opts[i](&params)
	}

	return &Background[T, R]{
		coworker:  num,
		golbalCtx: ctx,
		worker:    task,
		queux:     queux.NewQueuxMemory[*MessageWithResponce[T, R]](),
		params:    params,
	}
}
