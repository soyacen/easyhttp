package easyhttpretry

import (
	"net/http"
	"time"
)

var (
	DefaultRetriableStatusCode = []int{http.StatusRequestTimeout, http.StatusGatewayTimeout, http.StatusServiceUnavailable}
)

// Disable disables the retry behaviour on this call, or this interceptor.
//
// Its semantically the same to `WithMaxAttempts`
func Disable() Option {
	return WithMaxAttempts(0)
}

// WithMaxAttempts sets the maximum number of retries on this call, or this interceptor.
func WithMaxAttempts(maxAttempts uint) Option {
	return Option{applyFunc: func(o *options) {
		o.maxAttempts = maxAttempts
	}}
}

// WithBackoff sets the `BackoffFunc` used to control time between retries.
func WithBackoff(bf BackoffFunc) Option {
	return Option{applyFunc: func(o *options) {
		o.backoffFunc = bf
	}}
}

// WithStatusCodes sets which statusCodes should be retried.
func WithStatusCodes(codes ...int) Option {
	return Option{applyFunc: func(o *options) {
		o.statusCodes = codes
	}}
}

// WithPerRetryTimeout sets the timeout of each HTTP request
func WithPerRetryTimeout(timeout time.Duration) Option {
	return Option{applyFunc: func(o *options) {
		o.perCallTimeout = timeout
	}}
}

type options struct {
	maxAttempts    uint
	perCallTimeout time.Duration
	includeHeader  bool
	statusCodes    []int
	backoffFunc    BackoffFunc
}

type Option struct {
	applyFunc func(opt *options)
}

func defaultOptions() *options {
	o := &options{
		maxAttempts:    0,
		perCallTimeout: 0,
		includeHeader:  true,
		statusCodes:    DefaultRetriableStatusCode,
		backoffFunc: BackoffFunc(func(attempt uint) time.Duration {
			return BackoffLinearWithJitter(50*time.Millisecond, 0.10)(attempt)
		}),
	}
	return o
}

func (opt *options) apply(callOptions ...Option) *options {
	if len(callOptions) == 0 {
		return opt
	}
	for _, f := range callOptions {
		f.applyFunc(opt)
	}
	return opt
}
