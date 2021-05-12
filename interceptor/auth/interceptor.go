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
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		rawRequest.SetBasicAuth(username, password)
		return do(cli, req)
	}
}

type APIKeyAddTo int

const (
	AddToHeader APIKeyAddTo = 1
	AddToQuery  APIKeyAddTo = 2
)

func APIKey(key string, value string, addTo APIKeyAddTo) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if addTo == AddToHeader {
			if rawRequest.Header == nil {
				rawRequest.Header = make(http.Header)
			}
			rawRequest.Header.Set(key, value)
		} else if addTo == AddToQuery {
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
	return APIKey("Authorization", "Bearer "+token, AddToHeader)
}

func CustomToken(scheme, token string) easyhttp.Interceptor {
	if stringutils.IsBlank(scheme) {
		return APIKey("Authorization", token, AddToHeader)
	}
	return APIKey("Authorization", scheme+" "+token, AddToHeader)
}
