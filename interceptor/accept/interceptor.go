package easyhttpaccept

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/soyacen/bytebufferpool"
	"github.com/soyacen/easyhttp"
	"github.com/soyacen/goutils/ioutils"
)

var (
	kContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	kAcceptKey          = http.CanonicalHeaderKey("Accept")
	kContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
)

const (
	JsonContentType = "application/json"
	XMLContentType  = "application/xml"
)

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

func JSON(obj interface{}, unmarshalFunc ...func(data []byte, v interface{}) error) easyhttp.Interceptor {
	unmarshal := json.Unmarshal
	if len(unmarshalFunc) > 0 {
		unmarshal = unmarshalFunc[0]
	}
	return Object(obj, JsonContentType, unmarshal)
}

func XML(obj interface{}, unmarshalFunc ...func(data []byte, v interface{}) error) easyhttp.Interceptor {
	unmarshal := xml.Unmarshal
	if len(unmarshalFunc) > 0 {
		unmarshal = unmarshalFunc[0]
	}
	return Object(obj, XMLContentType, unmarshal)
}

func Object(obj interface{}, contentType string, unmarshalFunc func(data []byte, v interface{}) error) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		if unmarshalFunc == nil {
			return reply, errors.New("unmarshal function is nil")
		}
		if obj == nil {
			return reply, errors.New("object is nil")
		}
		if req.RawRequest().Header == nil {
			req.RawRequest().Header = make(http.Header)
		}
		req.RawRequest().Header.Set(kAcceptKey, contentType)
		reply, err = do(cli, req)
		if err != nil {
			return reply, err
		}
		rawResponse := reply.RawResponse()
		if reply == nil || rawResponse == nil {
			return reply, err
		}
		if rawResponse.ContentLength == 0 {
			return reply, err
		}

		if rawResponse.Body != nil {
			defer ioutils.CloseThrowError(rawResponse.Body, &err)
		}

		var body io.Reader
		body = rawResponse.Body
		cek := rawResponse.Header.Get(kContentEncodingKey)
		if strings.EqualFold(cek, "gzip") {
			if _, ok := rawResponse.Body.(*gzip.Reader); !ok {
				gzipReader, err := gzip.NewReader(rawResponse.Body)
				if err != nil {
					return reply, err
				}
				defer ioutils.CloseThrowError(gzipReader, &err)
				body = gzipReader
			}
		}

		ct := rawResponse.Header.Get(kContentTypeKey)
		if ct != contentType {
			return reply, fmt.Errorf("expexted content-type is %s, but actual content-type is %s", contentType, ct)
		}

		data, err := ioutil.ReadAll(body)
		if err != nil {
			return reply, err
		}
		err = unmarshalFunc(data, obj)
		if err != nil {
			return reply, err
		}
		return reply, err
	}
}
