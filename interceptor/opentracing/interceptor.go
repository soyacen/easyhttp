package easyhttpopentracing

import (
	"net/url"
	"sync"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	"github.com/soyacen/easyhttp"
)

func Opentracing(tracer opentracing.Tracer, opts ...Option) easyhttp.Interceptor {
	o := &options{
		operationName: "HTTP Client",
		componentName: "github.com/soyacen/easyhttp",
		urlTagFunc:    func(u *url.URL) string { return u.String() },
	}
	for _, opt := range opts {
		opt(o)
	}
	once := sync.Once{}
	transport := &nethttp.Transport{}
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		once.Do(func() {
			rawClient := cli.RawClient()
			if _, ok := rawClient.Transport.(*nethttp.Transport); !ok {
				transport.RoundTripper = rawClient.Transport
				rawClient.Transport = &nethttp.Transport{
					RoundTripper: rawClient.Transport,
				}
				cli.SetRawClient(rawClient)
			}
		})
		traceReq, ht := nethttp.TraceRequest(
			tracer,
			req.RawRequest(),
			nethttp.OperationName(o.operationName),
			nethttp.ComponentName(o.componentName),
			nethttp.URLTagFunc(o.urlTagFunc),
		)
		defer ht.Finish()
		req.SetRawRequest(traceReq)
		return do(cli, req)
	}
}
