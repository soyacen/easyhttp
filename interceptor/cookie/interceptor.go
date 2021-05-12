package easyhttpcookie

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/soyacen/easyhttp"
	"golang.org/x/net/publicsuffix"
)

func Cookies(cookies ...*http.Cookie) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header == nil {
			rawRequest.Header = make(http.Header)
		}
		for _, cookie := range cookies {
			rawRequest.AddCookie(cookie)
		}
		return do(cli, req)
	}
}

func Set(key, value string) easyhttp.Interceptor {
	return Cookies(&http.Cookie{Name: key, Value: value})
}

func DelAll() easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		if rawRequest.Header != nil {
			rawRequest.Header.Del("Cookie")
		}
		return do(cli, req)
	}
}

// SetMap sets a map of cookies represented by key-value pair.
func SetMap(cookies map[string]string) easyhttp.Interceptor {
	cks := make([]*http.Cookie, 0, len(cookies))
	for k, v := range cookies {
		cks = append(cks, &http.Cookie{Name: k, Value: v})
	}
	return Cookies(cks...)
}

// Jar creates a cookie jar to store HTTP cookies when they are sent down.
func Jar() easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		rawClient := cli.RawClient()
		rawClient.Jar = jar
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}
