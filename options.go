package worker

import "time"

type Option func(*params)

func WithTimout(timeout time.Duration) Option {
	return func(p *params) {
		p.Timeout = timeout
	}
}

func OnError(callback func(error)) Option {
	return func(p *params) {
		p.OnError = callback
	}
}
