package easyhttpcompression

import (
	"net/http"

	"github.com/soyacen/easyhttp"
)

func Enable(enabled bool) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawClient := cli.RawClient()
		transport, ok := rawClient.Transport.(*http.Transport)
		if !ok {
			return do(cli, req)
		}
		transport.DisableCompression = enabled
		rawClient.Transport = transport
		cli.SetRawClient(rawClient)
		return do(cli, req)
	}
}
