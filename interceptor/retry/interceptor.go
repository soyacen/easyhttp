package easyhttpretry

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/soyacen/bytebufferpool"

	"github.com/soyacen/easyhttp"
)

func Interceptor(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		// if attempt is 0, just call do.
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
			rawRequest.Body = ioutil.NopCloser(bodyReader)
			rawRequest.ContentLength = int64(bodyBuf.Len())
			rawRequest.GetBody = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(bytes.NewBuffer(bodyBuf.Bytes())), nil
			}
		}

		for attempt := uint(0); attempt <= o.maxAttempts; attempt++ {
			newRequest := spawnRequest(rawRequest, o.timeout, attempt)
			req.SetRawRequest(newRequest)

			if bodyReader != nil {
				// Reset the body reader after the request since at this point it's already read
				// Note that it's safe to ignore the error here since the 0,0 position is always valid
				_, _ = bodyReader.Seek(0, 0)
			}

			// call do
			reply, err = do(cli, req)

			// if have already retried the maximum number, return result
			if attempt == o.maxAttempts {
				return reply, err
			}

			// if raw context is deadline or canceled, just return result
			if rawCtx.Err() != nil {
				return reply, err
			}

			// check error retry condition
			if err != nil {
				// if this error should not retry, return result
				if !o.shouldRetryWithError(err) {
					return reply, err
				}
				// wait time
				if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
					return reply, err
				}
				continue
			}

			// reply or raw response is nil, must be intercept, just return
			if reply == nil || reply.RawResponse() == nil {
				return reply, err
			}

			// check code retry condition
			if reply != nil && reply.RawResponse() != nil {
				// check status code, if code should not retry, return result
				if !o.shouldRetryWithStatusCode(reply.RawResponse().StatusCode) {
					return reply, err
				}
				// if code is need retry,continue
				if e := waitRetryBackoff(attempt, rawCtx, o); e != nil {
					return reply, err
				}
				continue
			}
		}
		return reply, err
	}
}

func waitRetryBackoff(attempt uint, parentCtx context.Context, callOpts *options) error {
	waitTime := callOpts.backoffFunc(parentCtx, attempt)
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

func spawnRequest(rawRequest *http.Request, timeout time.Duration, attempt uint) *http.Request {
	ctx := rawRequest.Context()
	if timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}
	if attempt > 0 {
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		rawRequest.Header.Set("X-Retry", strconv.Itoa(int(attempt)))
	}
	return rawRequest.WithContext(ctx)
}
