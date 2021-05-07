package easyhttpopentracing

import (
	"sync"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/soyacen/easyhttp"
)

func Opentracing() easyhttp.Interceptor {
	once := sync.Once{}
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		once.Do(func() {
			cli.RawClient().Transport = &nethttp.Transport{
				RoundTripper: cli.RawClient().Transport,
			}
		})
		reqWithTrace, ht := nethttp.TraceRequest(tracer, req.RawRequest(), opts...)
		defer ht.Finish()
		newReq := req.WithRawRequest(reqWithTrace)

	}
}