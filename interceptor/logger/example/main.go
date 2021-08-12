package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/soyacen/easyhttp"
	easyhttplogger "github.com/soyacen/easyhttp/interceptor/logger"
)

func main() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttplogger.Interceptor(func(fields *easyhttplogger.Fields, reply *easyhttp.Reply) {
			log.Println(fields)
		})))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}
