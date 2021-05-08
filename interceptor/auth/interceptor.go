package easyhttpauth

import (
	"net/http"
	"net/url"

	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyhttp"
)

func BasicAuth(username, password string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		rawRequest.SetBasicAuth(username, password)
		return do(cli, req)
	}
}

type APIKeyAddTo int

const (
	APIKeyHeader APIKeyAddTo = 1
	APIKeyQuery  APIKeyAddTo = 2
)

func APIKey(key string, value string, addTo APIKeyAddTo) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if addTo == APIKeyHeader {
			if rawRequest.Header == nil {
				rawRequest.Header = make(http.Header)
			}
			rawRequest.Header.Set(key, value)
		} else if addTo == APIKeyQuery {
			query := make(url.Values)
			query.Set(key, value)
			reqURL := rawRequest.URL
			if stringutils.IsBlank(reqURL.RawQuery) {
				reqURL.RawQuery = query.Encode()
				reqURL.Query()
			} else {
				reqURL.RawQuery = reqURL.RawQuery + "&" + query.Encode()
			}
		}
		return do(cli, req)
	}
}

func BearerToken(token string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		rawRequest.Header.Set("Authorization", "Bearer "+token)
		return do(cli, req)
	}
}
