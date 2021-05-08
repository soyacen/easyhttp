package easyhttp

import (
	"context"
	"net/http"
)

type Request struct {
	rawRequest *http.Request
	opts       *executeOptions
}

func (r *Request) Context() context.Context {
	return r.rawRequest.Context()
}

func (r *Request) SetContext(ctx context.Context) {
	newRaw := r.rawRequest.WithContext(ctx)
	r.rawRequest = newRaw
}

func (r *Request) SetRawRequest(rawRequest *http.Request) {
	r.rawRequest = rawRequest
}

func (r *Request) RawRequest() *http.Request {
	return r.rawRequest
}
