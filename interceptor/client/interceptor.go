package easyhttpclient

import (
	"net/http"
	"time"

	"github.com/soyacen/easyhttp"
)

func Client(httpClient *http.Client) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		if httpClient != nil {
			cli.SetRawClient(httpClient)
		}
		return do(cli, req)
	}
}

func Timeout(timeout time.Duration) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		rawClient.Timeout = timeout
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}

func Jar(jar http.CookieJar) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		rawClient.Jar = jar
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}

func RedirectPolicy(redirectPolicy func(req *http.Request, via []*http.Request) error) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		rawClient.CheckRedirect = redirectPolicy
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}

func Transport(transport http.RoundTripper) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		rawClient.Transport = transport
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}
