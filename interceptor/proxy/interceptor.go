package easyhttpproxy

import (
	"net/http"
	"net/url"

	"github.com/soyacen/easyhttp"
)

// Set defines the proxy servers to be used based on the transport scheme
func Set(servers map[string]string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		transport, ok := rawClient.Transport.(*http.Transport)
		if !ok {
			return do(cli, req)
		}
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			if value, ok := servers[req.URL.Scheme]; ok {
				return url.Parse(value)
			}
			return http.ProxyFromEnvironment(req)
		}
		rawClient.Transport = transport
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}
