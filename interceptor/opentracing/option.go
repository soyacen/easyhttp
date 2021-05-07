package easyhttpopentracing

import (
	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

type Option = nethttp.ClientOption

// OperationName returns a ClientOption that sets the operation
// name for the client-side span.
var OperationName = nethttp.OperationName

// URLTagFunc returns a ClientOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
var URLTagFunc = nethttp.URLTagFunc

// ComponentName returns a ClientOption that sets the component
// name for the client-side span.
var ComponentName = nethttp.ComponentName

// ClientTrace returns a ClientOption that turns on or off
// extra instrumentation via httptrace.WithClientTrace.
var ClientTrace = nethttp.ClientTrace

// InjectSpanContext returns a ClientOption that turns on or off
// injection of the Span context in the request HTTP headers.
// If this option is not used, the default behaviour is to
// inject the span context.
var InjectSpanContext = nethttp.InjectSpanContext

// ClientSpanObserver returns a ClientOption that observes the span
// for the client-side span.
var ClientSpanObserver = nethttp.ClientSpanObserver
