package easyhttpurl

import (
	"net/url"
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

// PathParam replaces one or multiple path param expressions by the given value
func PathParam(key, value string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		rawRequest.URL.Path = replace(rawRequest.URL.Path, key, value)
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

// PathParams replaces one or multiple path param expressions by the given map of key-value pairs
func PathParams(params map[string]string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		for key, value := range params {
			rawRequest.URL.Path = replace(rawRequest.URL.Path, key, value)
		}
		req.SetRawRequest(rawRequest)
		return do(cli, req)
	}
}

func QueryParam(key string, values ...string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		query := make(url.Values)
		for _, iv := range values {
			query.Add(key, iv)
		}
		rawRequest := req.RawRequest()
		reqURL := rawRequest.URL
		if stringutils.IsBlank(reqURL.RawQuery) {
			reqURL.RawQuery = query.Encode()
		} else {
			reqURL.RawQuery = reqURL.RawQuery + "&" + query.Encode()
		}
		return do(cli, req)
	}
}

func QueryParams(param map[string][]string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		query := make(url.Values)
		for k, v := range param {
			for _, iv := range v {
				query.Add(k, iv)
			}
		}
		rawRequest := req.RawRequest()
		reqURL := rawRequest.URL
		if stringutils.IsBlank(reqURL.RawQuery) {
			reqURL.RawQuery = query.Encode()
		} else {
			reqURL.RawQuery = reqURL.RawQuery + "&" + query.Encode()
		}
		return do(cli, req)
	}
}

func replace(str, key, value string) string {
	return strings.Replace(str, ":"+key, value, -1)
}
