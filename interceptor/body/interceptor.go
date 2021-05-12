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
	kContentTypeKey   = http.CanonicalHeaderKey("Content-Type")
	kContentLengthKey = http.CanonicalHeaderKey("Content-Length")
)

const (
	kFormContentType = "application/x-www-form-urlencoded"
	JsonContentType  = "application/json"
	XMLContentType   = "application/xml"
	kTextContentType = "text/plain"
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
	rawRequest.Header.Set(kContentTypeKey, ct)
	rawRequest.Header.Set(kContentLengthKey, strconv.Itoa(bodyBuf.Len()))
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

func Reader(body io.Reader, contentType string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		bodyBuf := bytebufferpool.Get()
		bodyBuf.ReadFrom(body)
		defer bodyBuf.Free()
		setContent(bodyBuf, req, contentType)
		return do(cli, req)
	}
}

func Text(text string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		bodyBuf := bytebufferpool.Get()
		bodyBuf.WriteString(text)
		defer bodyBuf.Free()
		setContent(bodyBuf, req, kTextContentType)
		return do(cli, req)
	}
}

func Form(form url.Values) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		data := []byte(form.Encode())
		bodyBuf := bytebufferpool.Get()
		bodyBuf.Write(data)
		defer bodyBuf.Free()
		setContent(bodyBuf, req, kFormContentType)
		return do(cli, req)
	}
}

func JSON(obj interface{}, marshalFunc ...func(v interface{}) ([]byte, error)) easyhttp.Interceptor {
	marshal := json.Marshal
	if len(marshalFunc) > 0 {
		marshal = marshalFunc[0]
	}
	return Object(obj, JsonContentType, marshal)
}

func XML(obj interface{}, marshalFunc ...func(v interface{}) ([]byte, error)) easyhttp.Interceptor {
	marshal := xml.Marshal
	if len(marshalFunc) > 0 {
		marshal = marshalFunc[0]
	}
	return Object(obj, XMLContentType, marshal)
}

func Object(obj interface{}, contentType string, marshalFunc func(v interface{}) ([]byte, error)) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		bodyBuf := bytebufferpool.Get()
		defer bodyBuf.Free()
		if err := writeObj(obj, bodyBuf, marshalFunc); err != nil {
			return nil, err
		}
		setContent(bodyBuf, req, contentType)
		return do(cli, req)
	}
}
