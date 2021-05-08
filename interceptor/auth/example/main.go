package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/soyacen/easyhttp"
	easyhttpauth "github.com/soyacen/easyhttp/interceptor/auth"
)

func ExampleBearerToken() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttpauth.BearerToken("dsade1232few23r2e")))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func ExampleBasicAuth() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttpauth.BasicAuth("fakeuser", "fakepwd")))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func ExampleAPIKey() {
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(easyhttpauth.APIKey("token", "faketoken", easyhttpauth.APIKeyQuery)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}
