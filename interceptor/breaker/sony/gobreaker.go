package easyhttpsonybreaker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sony/gobreaker"

	"github.com/soyacen/easyhttp"
)

func Interceptor(Name string, opts ...Option) easyhttp.Interceptor {
	st := defaultSettings(Name)
	apply(st, opts...)
	cb := gobreaker.NewCircuitBreaker(*st)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		result, err := cb.Execute(func() (interface{}, error) {
			reply, err = do(cli, req)
			if err != nil {
				return nil, err
			}
			if reply == nil {
				return nil, errors.New("reply is nil")
			}
			if reply.RawResponse() == nil {
				return nil, errors.New("http response is nil")
			}
			if reply.RawResponse().StatusCode >= http.StatusInternalServerError {
				return nil, fmt.Errorf("server returned %d status code", reply.RawResponse().StatusCode)
			}
			return reply, nil
		})
		if result != nil {
			reply = result.(*easyhttp.Reply)
		}
		return reply, err
	}
}
