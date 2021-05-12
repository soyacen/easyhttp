package easyhttplogging

import (
	"net/http"
	"time"

	"github.com/soyacen/easyhttp"
)

type Fields struct {
	System         string
	Method         string
	URL            string
	StartTime      time.Time
	Deadline       time.Time
	RequestHeader  http.Header
	Error          error
	requestBody    []byte
	ResponseHeader http.Header
	Status         string
	StatusCode     int
}

func Logger(logFunc func(fields *Fields, reply *easyhttp.Reply)) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		startTime := time.Now()
		fields := &Fields{
			System:    "http.client",
			StartTime: startTime,
		}
		if d, ok := req.RawRequest().Context().Deadline(); ok {
			fields.Deadline = d
		}
		rawRequest := req.RawRequest()
		fields.Method = rawRequest.Method
		fields.URL = rawRequest.URL.String()
		fields.RequestHeader = rawRequest.Header.Clone()
		r, err := do(cli, req)
		rawResponse := r.RawResponse()
		fields.Status = rawResponse.Status
		fields.StatusCode = rawResponse.StatusCode
		fields.ResponseHeader = rawResponse.Header.Clone()
		fields.Error = err
		logFunc(fields, reply)
		return r, err
	}
}
