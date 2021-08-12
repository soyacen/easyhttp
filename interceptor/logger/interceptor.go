package easyhttplogger

import (
	"time"

	"github.com/soyacen/easyhttp"
)

func Interceptor(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		if o.loggerFactory == nil {
			return do(cli, req)
		}
		logger := o.loggerFactory(req.RawRequest().Context())
		startTime := time.Now()
		r, err := do(cli, req)
		rawRequest := req.RawRequest()
		rawResponse := r.RawResponse()
		builder := NewFieldBuilder().
			System().
			StartTime(startTime).
			Deadline(req.RawRequest().Context()).
			Method(rawRequest.Method).
			URI(rawRequest.URL.String()).
			RequestHeader(rawRequest.Header.Clone()).
			Latency(time.Since(startTime)).
			Status(rawResponse.Status).
			StatusCode(rawResponse.StatusCode).
			ResponseHeader(rawResponse.Header.Clone()).
			Error(err)
		logger.Log(builder.Build())
		return r, err
	}
}
