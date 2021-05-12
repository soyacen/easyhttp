package easyhttptls

import (
	"crypto/tls"
	"net/http"

	"github.com/soyacen/easyhttp"
)

// Config defines the request TLS connection config
func Config(config *tls.Config) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		transport, ok := rawClient.Transport.(*http.Transport)
		if !ok {
			return do(cli, req)
		}
		transport.TLSClientConfig = config
		rawClient.Transport = transport
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}
