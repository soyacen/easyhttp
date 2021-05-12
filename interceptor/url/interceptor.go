package easyhttpurl

import (
	"strings"

	"github.com/soyacen/easyhttp"
	"github.com/soyacen/goutils/stringutils"
)

// Schema parses and defines a schema values in the outgoing request
func Schema(schema string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if stringutils.IsBlank(rawRequest.URL.Scheme) {
			rawRequest.URL.Scheme = schema
		}
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

// Host parses and defines a host URL values in the outgoing request
func Host(host string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if stringutils.IsBlank(rawRequest.URL.Host) {
			rawRequest.URL.Host = host
		}
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

// Param replaces one or multiple path param expressions by the given value
func Param(key, value string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		rawRequest.URL.Path = replace(rawRequest.URL.Path, key, value)
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

// Params replaces one or multiple path param expressions by the given map of key-value pairs
func Params(params map[string]string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		for key, value := range params {
			rawRequest.URL.Path = replace(rawRequest.URL.Path, key, value)
		}
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

func replace(str, key, value string) string {
	return strings.Replace(str, ":"+key, value, -1)
}
