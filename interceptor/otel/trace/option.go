package oteltrace

import (
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type SpanNameFunc func(u *url.URL) string

type options struct {
	propagator   propagation.TextMapPropagator
	tracer       trace.Tracer
	spanNameFunc SpanNameFunc
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(o *options)

func defaultOptions() *options {
	return &options{
		tracer:       otel.Tracer(""),
		propagator:   otel.GetTextMapPropagator(),
		spanNameFunc: func(u *url.URL) string { return u.EscapedPath() },
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return func(o *options) {
		o.tracer = tracer
	}
}

func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(o *options) {
		o.propagator = propagator
	}
}

func WithSpanNameFunc(f func(u *url.URL) string) Option {
	return func(options *options) {
		options.spanNameFunc = f
	}
}
