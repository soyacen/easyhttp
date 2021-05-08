package easyhttpopentracing

import (
	"sync"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	"github.com/soyacen/easyhttp"
)

func Opentracing(tracer opentracing.Tracer, opts ...Option) easyhttp.Interceptor {
	once := sync.Once{}
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		once.Do(func() {
			cli.RawClient().Transport = &nethttp.Transport{
				RoundTripper: cli.RawClient().Transport,
			}
		})
		traceReq, ht := nethttp.TraceRequest(tracer, req.RawRequest(), opts...)
		defer ht.Finish()
		req.SetRawRequest(traceReq)
		return do(cli, req)
	}
}
