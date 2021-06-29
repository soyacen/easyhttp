package easyhttpopentracing

import (
	"net/url"
)

type options struct {
	operationName string
	componentName string
	urlTagFunc    func(u *url.URL) string
}

type Option func(o *options)

// OperationName returns a ClientOption that sets the operation
// name for the client-side span.
func OperationName(operationName string) Option {
	return func(o *options) {
		o.operationName = operationName
	}
}

// URLTagFunc returns a ClientOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
func URLTagFunc(f func(u *url.URL) string) Option {
	return func(options *options) {
		options.urlTagFunc = f
	}
}

// ComponentName returns a ClientOption that sets the component
// name for the client-side span.
func ComponentName(componentName string) Option {
	return func(options *options) {
		options.componentName = componentName
	}
}
