package easyhttpretry

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/soyacen/bytebufferpool"
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

		var bodyReader *bytes.Reader
		bodyBuf := bytebufferpool.Get()
		defer bodyBuf.Free()
		if rawRequest.Body != nil {
			_, err := bodyBuf.ReadFrom(rawRequest.Body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(bodyBuf.Bytes())
			rawRequest.Body = io.NopCloser(bodyReader)
			rawRequest.ContentLength = int64(bodyBuf.Len())
			buf := bodyBuf.Bytes()
			rawRequest.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(buf)), nil
			}
		}

		for attempt := uint(0); attempt <= o.maxAttempts; attempt++ {
			rawRequest := requestWithPerTimeout(rawRequest, o, attempt)
			req.SetRawRequest(rawRequest)

			if bodyReader != nil {
				// Reset the body reader after the request since at this point it's already read
				// Note that it's safe to ignore the error here since the 0,0 position is always valid
				_, _ = bodyReader.Seek(0, 0)
			}
			reply, err = do(cli, req)
			if attempt == o.maxAttempts {
				return reply, err
			}
			if err != nil {
				if !o.shouldRetryWithError(err) {
					return reply, err
				}
				// if raw context is deadline or canceled, just return
				if rawCtx.Err() != nil {
					return reply, err
				}
				if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
					return reply, err
				}
				continue
			}
			// reply or raw response is nil, must be intercept, just return
			if reply == nil || reply.RawResponse() == nil {
				return reply, err
			}

			// check status code
			if !o.shouldRetryWithStatusCode(reply.RawResponse().StatusCode) {
				return reply, err
			}
			// if raw context is deadline or canceled, just return
			if rawCtx.Err() != nil {
				return reply, err
			}
			// if code is need retry,continue
			if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
				return reply, err
			}
			continue
		}
		return reply, err
	}
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
	if callOpts.timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, callOpts.timeout)
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
