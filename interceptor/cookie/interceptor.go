package easyhttpcookie

import (
	"net/http"

	"github.com/soyacen/easyhttp"
)

func Cookies(cookies ...*http.Cookie) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		rawRequest := req.RawRequest()
		for _, cookie := range cookies {
			rawRequest.AddCookie(cookie)
		}
		return do(cli, req)
	}
}
