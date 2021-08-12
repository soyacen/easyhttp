package easyhttpmultipart

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/soyacen/bytebufferpool"
	"github.com/soyacen/goutils/ioutils"
	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyhttp"
)

type FormData struct {
	// fieldName is form field name
	fieldName string
	// data is text form
	data string
	// filepath file form,
	filepath string
	// fileContent is happy as long as you pass it a io.ReadCloser (which most file use anyways)
	fileContent io.ReadCloser
	// fileMime represents which mimetime should be sent along with the file.
	// When empty, defaults to application/octet-stream
	fileMime string
}

func NewFile(fieldName string, filepath string, fileContent io.ReadCloser, fileMime string) *FormData {
	return &FormData{
		fieldName:   fieldName,
		filepath:    filepath,
		fileContent: fileContent,
		fileMime:    fileMime,
	}
}

func NewData(fieldName string, data string) *FormData {
	return &FormData{
		fieldName: fieldName,
		data:      data,
	}
}

func Interceptor(formData ...*FormData) easyhttp.Interceptor {
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		for _, file := range formData {
			if stringutils.IsNotEmpty(file.data) {
				continue
			}
			if file.fileContent != nil {
				continue
			}
			// open file by filepath
			if stringutils.IsNotBlank(file.filepath) {
				fd, err := os.Open(file.filepath)
				if err != nil {
					return nil, err
				}
				file.fileContent = fd
			}
		}

		requestBody := bytebufferpool.Get()

		multipartWriter := multipart.NewWriter(requestBody)
		defer ioutils.CloseThrowError(multipartWriter, &err)

		for i, f := range formData {
			fieldName := f.fieldName
			if f.fileContent != nil {
				// generate default field name
				if stringutils.IsBlank(fieldName) {
					if len(formData) > 1 {
						fieldName = strings.Join([]string{"file", strconv.Itoa(i + 1)}, "")
					} else {
						fieldName = "file"
					}
				}
				// write file
				err = writeFile(fieldName, f, multipartWriter)
				if err != nil {
					return nil, err
				}
			} else {
				// write data
				err = multipartWriter.WriteField(fieldName, f.data)
				if err != nil {
					return nil, err
				}
			}

		}
		if req.RawRequest().Header == nil {
			req.RawRequest().Header = make(http.Header)
		}
		req.RawRequest().Header.Add("Content-Type", multipartWriter.FormDataContentType())

		return do(cli, req)
	}
}

func writeFile(fieldName string, f *FormData, multipartWriter *multipart.Writer) (err error) {
	var filename string
	if stringutils.IsNotBlank(f.filepath) {
		filename = filepath.Base(filename)
	} else {
		filename = "filename"
	}

	var writer io.Writer

	if stringutils.IsNotBlank(f.fileMime) {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(filename)))
		h.Set("Content-Type", f.fileMime)
		writer, err = multipartWriter.CreatePart(h)
	} else {
		writer, err = multipartWriter.CreateFormFile(fieldName, filename)
	}

	if err != nil {
		return err
	}

	if _, err = io.Copy(writer, f.fileContent); err != nil && err != io.EOF {
		return err
	}
	return f.fileContent.Close()
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
