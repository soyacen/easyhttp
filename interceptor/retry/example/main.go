package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"
	"time"

	"github.com/soyacen/easyhttp"

	easyhttpretry "github.com/soyacen/easyhttp/interceptor/retry"
)

//func Example500() {
func main() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/status/500",
		easyhttp.ChainInterceptor(
			easyhttpretry.Interceptor(
				easyhttpretry.WithBackoff(easyhttpretry.BackoffLinear(time.Second)),
				easyhttpretry.WithTimeout(time.Second),
				easyhttpretry.WithMaxAttempts(2),
			)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func Example400() {
	//func main(){
	client := easyhttp.NewClient()
	timeout, _ := context.WithTimeout(context.Background(), time.Hour)
	reply, err := client.Get(
		timeout,
		"http://httpbin.org/status/400",
		easyhttp.ChainInterceptor(
			easyhttpretry.Interceptor(
				easyhttpretry.WithTimeout(time.Second),
				easyhttpretry.WithMaxAttempts(2),
			)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func Example200() {
	//func main() {
	client := easyhttp.NewClient()
	timeout, _ := context.WithTimeout(context.Background(), time.Hour)
	reply, err := client.Get(
		timeout,
		"http://httpbin.org/status/200",
		easyhttp.ChainInterceptor(
			easyhttpretry.Interceptor(
				easyhttpretry.WithBackoff(easyhttpretry.BackoffLinear(time.Second)),
				easyhttpretry.WithTimeout(time.Second),
				easyhttpretry.WithMaxAttempts(2),
			)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}
