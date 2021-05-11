package easyhttpretry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/soyacen/easyhttp"
)

func Retry(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		// maxAttempts = 0, don't need retry, short circuit for simplicity, and avoiding allocations.
		if o.maxAttempts == 0 {
			return do(cli, req)
		}
		rawRequest := req.RawRequest()
		rawCtx := rawRequest.Context()
		for attempt := uint(0); attempt <= o.maxAttempts; attempt++ {
			rawRequest := requestWithPerTimeout(rawRequest, o, attempt)
			req.SetRawRequest(rawRequest)
			fmt.Println("call ...")
			reply, err = do(cli, req)

			// if check error is not nil
			if err != nil {
				// if raw context is deadline or canceled, just return
				if rawCtx.Err() != nil {
					return reply, err
				}
				//  when err is context.DeadlineExceeded or context.Canceled, retry again
				if isDeadlineOrCanceled(err) {
					if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
						return reply, err
					}
					continue
				} else {
					// other error, just return
					return reply, err
				}
			}

			// reply must be intercept, just return
			if reply == nil || reply.RawResponse() == nil {
				return reply, err
			}

			// check status code
			if checkStatusCode(reply.RawResponse().StatusCode, o.statusCodes) {
				// if raw context is deadline or canceled, just return
				if rawCtx.Err() != nil {
					return
				}
				// if code is need retry,continue
				if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
					return reply, err
				}
				continue
			}

			// other, return
			return reply, err
		}
		return
	}
}

func checkStatusCode(statusCode int, codes []int) bool {
	for _, code := range codes {
		if code == statusCode {
			return true
		}
	}
	return false
}

func waitRetryBackoff(attempt uint, parentCtx context.Context, callOpts *options) error {
	waitTime := callOpts.backoffFunc(attempt)
	fmt.Println(waitTime)
	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		select {
		case <-parentCtx.Done():
			timer.Stop()
			return parentCtx.Err()
		case <-timer.C:
		}
	}
	return nil
}

func requestWithPerTimeout(rawRequest *http.Request, callOpts *options, attempt uint) *http.Request {
	ctx := rawRequest.Context()
	if callOpts.perCallTimeout != 0 {
		ctx, _ = context.WithTimeout(ctx, callOpts.perCallTimeout)
	}
	if attempt > 0 && callOpts.includeHeader {
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		rawRequest.Header.Set("x-retry-attempt", strconv.Itoa(int(attempt)))
	}
	return rawRequest.WithContext(ctx)
}

func isDeadlineOrCanceled(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
