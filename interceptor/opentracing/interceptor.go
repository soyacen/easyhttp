package easyhttpopentracing

import (
	"net/http"
	"net/url"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	"github.com/soyacen/easyhttp"
)

func Transport(rawTransport http.RoundTripper) http.RoundTripper {
	return &nethttp.Transport{
		RoundTripper: rawTransport,
	}
}

func Interceptor(tracer opentracing.Tracer, opts ...Option) easyhttp.Interceptor {
	o := &options{
		operationName: "HTTP Client",
		componentName: "github.com/soyacen/easyhttp",
		urlTagFunc:    func(u *url.URL) string { return u.String() },
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
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
