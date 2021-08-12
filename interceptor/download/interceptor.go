package easyhttpdownload

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	filepathutils "path/filepath"
	"strings"

	"github.com/soyacen/goutils/ioutils"

	"github.com/soyacen/easyhttp"
)

var (
	kContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	kAcceptKey          = http.CanonicalHeaderKey("Accept")
	kContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
)

func Interceptor(filepath string) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		err = createDirectory(filepathutils.Dir(filepath))
		if err != nil {
			return nil, err
		}
		if _, err = os.Stat(filepath); err != nil {
			if os.IsExist(err) {
				return
			}
		}

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

		fd, err := os.Create(filepath)
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

func createDirectory(dir string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return
			}
		}
	}
	return
}
