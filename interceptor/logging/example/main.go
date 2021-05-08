package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/soyacen/easyhttp"
	easyhttplogging "github.com/soyacen/easyhttp/interceptor/logging"
)

func main() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttplogging.Logger(func(fields *easyhttplogging.Fields) {
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
