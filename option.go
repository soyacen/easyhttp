package easyhttp

import (
	"net/http"
)

type executeOptions struct {
	interceptors []Interceptor
	interceptor  Interceptor

	client *http.Client
}

func defaultExecuteOptions() *executeOptions {
	o := &executeOptions{
		interceptors: make([]Interceptor, 0),
	}
	return o
}

func (o *executeOptions) apply(opts ...ExecuteOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type ExecuteOption func(o *executeOptions)

func ChainInterceptor(interceptors ...Interceptor) ExecuteOption {
	return func(o *executeOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}
