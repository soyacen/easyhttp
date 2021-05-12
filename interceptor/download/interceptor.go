package easyhttpdownload

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/soyacen/easyhttp"
	"github.com/soyacen/goutils/ioutils"
)

var (
	kContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	kAcceptKey          = http.CanonicalHeaderKey("Accept")
	kContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
)

func Download(filename string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
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

		fd, err := os.Create(filename)
		if err != nil {
			return reply, err
		}
		defer ioutils.CloseThrowError(fd, &err)

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

		if _, err := io.Copy(fd, body); err != nil && err != io.EOF {
			return reply, err
		}
		return reply, err
	}
}
