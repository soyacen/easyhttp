package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/soyacen/easyhttp"

	easyhttpheader "github.com/soyacen/easyhttp/interceptor/header"
)

func main() {
	client := easyhttp.NewClient(easyhttp.WithChainInterceptor(easyhttpheader.SetHeader("Authorization", "abcdefghijklmnopqrstuvwxyz")))
	reply1, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttpheader.SetHeader("X-Forward-For", "127.93.4.5")))
	if err != nil {
		log.Fatalln(err)
	}
	response1, err := httputil.DumpResponse(reply1.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response1))

	reply2, err := client.Post(
		context.Background(),
		"http://httpbin.org/post",
		easyhttp.ChainInterceptor(easyhttpheader.SetHeader("X-Trace-ID", "1id9dj1jdo0")))
	if err != nil {
		log.Fatalln(err)
	}
	response2, err := httputil.DumpResponse(reply2.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response2))
}
