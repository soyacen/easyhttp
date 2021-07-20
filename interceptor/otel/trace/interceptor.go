package easyhttpopentracing

import (
	"io"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"

	"github.com/soyacen/easyhttp"
)

func Interceptor(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		ctx, span := o.tracer.Start(
			req.Context(),
			o.spanNameFunc(req.RawRequest().URL),
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				semconv.RPCSystemKey.String("http"),
				semconv.HTTPMethodKey.String(req.RawRequest().Method),
				semconv.HTTPURLKey.String(req.RawRequest().URL.String()),
				semconv.HTTPFlavorKey.String(req.RawRequest().Proto),
			),
		)
		o.propagator.Inject(ctx, propagation.HeaderCarrier(req.RawRequest().Header))
		req.SetContext(ctx)
		reply, err = do(cli, req)
		if err != nil && err != io.EOF {
			span.RecordError(err)
			span.SetStatus(codes.Error, status.Code(err).String())
		} else {
			span.SetStatus(codes.Ok, "OK")
		}
		if reply != nil && reply.RawResponse() != nil {
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(reply.RawResponse().StatusCode),
			)
		}
		span.End()
		return reply, err
	}
}
