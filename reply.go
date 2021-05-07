package easyhttp

import "net/http"

type Reply struct {
	rawResponse *http.Response
	rawRequest  *http.Request
}

func (r *Reply) RawResponse() *http.Response {
	return r.rawResponse
}

func (r *Reply) RawRequest() *http.Request {
	return r.rawRequest
}
