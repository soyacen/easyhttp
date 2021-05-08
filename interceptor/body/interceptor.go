package easyhttpbody

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/soyacen/bytebufferpool"

	"github.com/soyacen/easyhttp"
)

var (
	ContentTypeKey   = http.CanonicalHeaderKey("Content-Type")
	ContentLengthKey = http.CanonicalHeaderKey("Content-Length")
)

const (
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
)

func setContent(bodyBuf *bytebufferpool.ByteBuffer, req *easyhttp.Request, ct string) {
	rawRequest := req.RawRequest()
	rawRequest.Body = io.NopCloser(bytes.NewReader(bodyBuf.Bytes()))
	rawRequest.ContentLength = int64(bodyBuf.Len())
	buf := bodyBuf.Bytes()
	rawRequest.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(buf)), nil
	}
	if rawRequest.Header == nil {
		rawRequest.Header = make(http.Header)
	}
	rawRequest.Header.Set(ContentTypeKey, ct)
	rawRequest.Header.Set(ContentLengthKey, strconv.Itoa(bodyBuf.Len()))
}

func writeObj(obj interface{}, bodyBuf *bytebufferpool.ByteBuffer, marshal func(v interface{}) ([]byte, error)) error {
	switch obj.(type) {
	case string:
		bodyBuf.WriteString(obj.(string))
	case []byte:
		bodyBuf.Write(obj.([]byte))
	default:
		data, err := marshal(obj)
		if err != nil {
			return err
		}
		bodyBuf.Write(data)
	}
	return nil
}

func Form(form url.Values) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		formData := make(url.Values)
		for k, v := range form {
			for _, iv := range v {
				formData.Add(k, iv)
			}
		}
		data := []byte(formData.Encode())
		bodyBuf := bytebufferpool.Get()
		bodyBuf.Write(data)
		defer bodyBuf.Free()
		setContent(bodyBuf, req, FormContentType)
		return do(cli, req)
	}
}

func JSON(obj interface{}, marshalFunc ...func(v interface{}) ([]byte, error)) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		marshal := json.Marshal
		if len(marshalFunc) > 0 {
			marshal = marshalFunc[0]
		}
		bodyBuf := bytebufferpool.Get()
		defer bodyBuf.Free()
		if err := writeObj(obj, bodyBuf, marshal); err != nil {
			return nil, err
		}
		setContent(bodyBuf, req, JsonContentType)
		return do(cli, req)
	}
}

func XML(obj interface{}, marshalFunc ...func(v interface{}) ([]byte, error)) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		marshal := xml.Marshal
		if len(marshalFunc) > 0 {
			marshal = marshalFunc[0]
		}
		bodyBuf := bytebufferpool.Get()
		defer bodyBuf.Free()
		if err := writeObj(obj, bodyBuf, marshal); err != nil {
			return nil, err
		}
		setContent(bodyBuf, req, JsonContentType)
		return do(cli, req)
	}
}
