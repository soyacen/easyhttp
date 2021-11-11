package easyhttpretry

import (
	"context"
	"crypto/x509"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/soyacen/goutils/backoffutils"
)

var (
	defaultRetriableStatusCode = []int{
		0, // means did not get a response. need to retry
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusLocked,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusInsufficientStorage,
	}

	defaultRetriableError = []error{
		context.DeadlineExceeded,
		context.Canceled,
	}
)

// RetryWithError whether the request should be retried based on error
type RetryWithError func(err error) bool

// RetryWithStatusCode  whether the request should be retried based on statusCode
type RetryWithStatusCode func(statusCode int) bool

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
func WithBackoff(bf backoffutils.BackoffFunc) Option {
	return Option{applyFunc: func(o *options) {
		o.backoffFunc = bf
	}}
}

// WithRetryPolicy sets retry policy which decides if a request should be retried.
func WithRetryPolicy(policyByError RetryWithError, policyByStatusCode RetryWithStatusCode) Option {
	return Option{applyFunc: func(o *options) {
		o.shouldRetryWithError = policyByError
		o.shouldRetryWithStatusCode = policyByStatusCode
	}}
}

// WithTimeout sets the timeout of each HTTP request
func WithTimeout(timeout time.Duration) Option {
	return Option{applyFunc: func(o *options) {
		o.timeout = timeout
	}}
}

type options struct {
	maxAttempts               uint
	timeout                   time.Duration
	backoffFunc               backoffutils.BackoffFunc
	shouldRetryWithError      RetryWithError
	shouldRetryWithStatusCode RetryWithStatusCode
}

type Option struct {
	applyFunc func(opt *options)
}

func defaultOptions() *options {
	o := &options{
		maxAttempts:               0,
		timeout:                   0,
		shouldRetryWithError:      defaultRetryWithError,
		shouldRetryWithStatusCode: defaultRetryWithStatusCode,
		backoffFunc:               BackoffExponentialWithJitter(50*time.Millisecond, 0.10),
	}
	return o
}

func (opt *options) apply(callOptions ...Option) {
	for _, f := range callOptions {
		f.applyFunc(opt)
	}
}

func defaultRetryWithError(err error) bool {
	// check if error is of type temporary
	t, ok := err.(interface{ Temporary() bool })
	if ok && t.Temporary() {
		return true
	}

	// we cannot know all errors, so we filter errors that should NOT be retried
	switch e := err.(type) {
	case *url.Error:
		switch {
		case
			e.Op == "parse",
			strings.Contains(e.Err.Error(), "stopped after"),
			strings.Contains(e.Error(), "unsupported protocol scheme"),
			strings.Contains(e.Error(), "no Host in request URL"):
			return false
		}
		// check inner error of url.Error
		switch e.Err.(type) {
		case // this errors will not likely change when retrying
			x509.UnknownAuthorityError,
			x509.CertificateInvalidError,
			x509.ConstraintViolationError:
			return false
		}
	}

	for _, e := range defaultRetriableError {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

func defaultRetryWithStatusCode(statusCode int) bool {
	for _, code := range defaultRetriableStatusCode {
		if code == statusCode {
			return true
		}
	}
	return false
}

// BackoffExponentialWithJitter creates an exponential backoff like
// BackoffExponential does, but adds jitter.
func BackoffExponentialWithJitter(scalar time.Duration, jitterFraction float64) backoffutils.BackoffFunc {
	return func(ctx context.Context, attempt uint) time.Duration {
		return backoffutils.JitterUp(backoffutils.Exponential(scalar), jitterFraction)(ctx, attempt)
	}
}
