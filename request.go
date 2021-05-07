package easyhttp

import "net/http"

type Request struct {
	rawRequest *http.Request
	opts       *executeOptions
}

func (r *Request) RawRequest() *http.Request {
	return r.rawRequest
}
