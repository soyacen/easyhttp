package easyhttpheader

import (
	"net/http"

	"github.com/soyacen/easyhttp"
)

func SetHeader(key, value string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		rawRequest.Header.Set(key, value)
		return do(cli, req)
	}
}

func SetHeaders(headers map[string]string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		for key, value := range headers {
			rawRequest.Header.Set(key, value)
		}
		return do(cli, req)
	}
}

func AddHeaders(key string, values ...string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		for _, value := range values {
			rawRequest.Header.Add(key, value)
		}
		return do(cli, req)
	}
}
