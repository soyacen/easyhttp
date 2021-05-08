package easyhttpquery

import (
	"net/url"

	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyhttp"
)

func SetQueryParams(param map[string][]string) easyhttp.Interceptor {
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
