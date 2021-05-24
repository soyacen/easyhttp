package easyhttpresp

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
	"github.com/soyacen/goutils/stringutils"
	"google.golang.org/protobuf/proto"
)

var (
	kContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	kAcceptKey          = http.CanonicalHeaderKey("Accept")
	kContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
)

const (
	kJsonContentType     = "application/json"
	kXMLContentType      = "application/xml"
	kProtobufContentType = "application/x-protobuf"
)

type options struct {
	contentType      string
	unmarshalFunc    func(data []byte, v interface{}) error
	checkContentType bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultOptions() *options {
	return &options{
		contentType:      "",
		unmarshalFunc:    nil,
		checkContentType: false,
	}
}

type Option func(o *options)

func Accept(contentType string) Option {
	return func(o *options) {
		o.contentType = contentType
	}
}

func UnmarshalFunc(unmarshalFunc func(data []byte, v interface{}) error) Option {
	return func(o *options) {
		o.unmarshalFunc = unmarshalFunc
	}
}

func CheckContentType() Option {
	return func(o *options) {
		o.checkContentType = true
	}
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

func JSON(obj interface{}, opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.contentType = kJsonContentType
	o.unmarshalFunc = json.Unmarshal
	o.apply(opts...)
	return object(obj, o)
}

func XML(obj interface{}, opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.contentType = kXMLContentType
	o.unmarshalFunc = xml.Unmarshal
	o.apply(opts...)
	return object(obj, o)
}

func Protobuf(obj proto.Message, opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.contentType = kXMLContentType
	o.unmarshalFunc = func(data []byte, v interface{}) error {
		message, ok := v.(proto.Message)
		if !ok {
			return errors.New("failed convert obj to proto.Message")
		}
		return proto.Unmarshal(data, message)
	}
	o.apply(opts...)
	return object(obj, o)
}

func Object(obj interface{}, opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	return object(obj, o)
}

func object(obj interface{}, o *options) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		if obj == nil {
			return reply, errors.New("object is nil")
		}
		if o.unmarshalFunc == nil {
			return reply, errors.New("unmarshal function is nil")
		}
		if req.RawRequest().Header == nil {
			req.RawRequest().Header = make(http.Header)
		}
		req.RawRequest().Header.Set(kAcceptKey, o.contentType)
		reply, err = do(cli, req)
		if err != nil {
			return reply, err
		}
		rawResponse := reply.RawResponse()
		if reply == nil || rawResponse == nil {
			return reply, err
		}
		if rawResponse.Body != nil {
			defer ioutils.CloseThrowError(rawResponse.Body, &err)
		}
		if rawResponse.ContentLength == 0 {
			return reply, err
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
		if o.checkContentType && stringutils.IsNotBlank(ct) && !strings.Contains(ct, o.contentType) {
			return reply, fmt.Errorf("expexted content-type is %s, but actual content-type is %s", o.contentType, ct)
		}

		data, err := ioutil.ReadAll(body)
		if err != nil {
			return reply, err
		}
		err = o.unmarshalFunc(data, obj)
		if err != nil {
			return reply, err
		}
		return reply, err
	}
}
