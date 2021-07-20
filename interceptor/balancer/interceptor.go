package easyhttpbalancer

import (
	"github.com/soyacen/easyhttp"
)

func Interceptor(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		pickerInfo := PickerInfo{
			URL:    req.RawRequest().URL,
			Header: req.RawRequest().Header,
			Ctx:    req.RawRequest().Context(),
		}
		pickResult, err := o.picker.Pick(pickerInfo)
		if err != nil {
			return
		}
		req.RawRequest().URL.Host = pickResult.Host
		reply, err = do(cli, req)
		return reply, err
	}
}
