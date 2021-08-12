package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/soyacen/easyhttp"
	easyhttplogger "github.com/soyacen/easyhttp/interceptor/logger"
)

type Logger struct {
}

func (l *Logger) Log(fields map[string]interface{}) {
	fmt.Println(fields)
}

func main() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttplogger.Interceptor(easyhttplogger.WithLoggerFactory(
			func(ctx context.Context) easyhttplogger.Logger {
				return &Logger{}
			}))))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}
