package easyhttpopentracing

import (
	"sync"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	"github.com/soyacen/easyhttp"
)

func Opentracing(tracer opentracing.Tracer, opts ...Option) easyhttp.Interceptor {
	once := sync.Once{}
	transport := &nethttp.Transport{}
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		once.Do(func() {
			rawClient := cli.RawClient()
			transport.RoundTripper = rawClient.Transport
			rawClient.Transport = &nethttp.Transport{
				RoundTripper: rawClient.Transport,
			}
			cli.SetRawClient(rawClient)
		})
		traceReq, ht := nethttp.TraceRequest(tracer, req.RawRequest(), opts...)
		defer ht.Finish()
		req.SetRawRequest(traceReq)
		return do(cli, req)
	}
}
