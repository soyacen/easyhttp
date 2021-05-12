package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/soyacen/easyhttp"
)

func main() {
	client := easyhttp.NewClient(
		easyhttp.WithChainInterceptor(
			func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
				log.Println("enter interceptor 1")
				defer log.Println("exit interceptor 1")
				return do(cli, req)
			},
			func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
				log.Println("enter interceptor 2")
				defer log.Println("exit interceptor 2")
				return do(cli, req)
			},
		),
	)
	response, err := client.Get(
		context.Background(),
		"http://httpbin.org/get",
		easyhttp.ChainInterceptor(
			func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
				log.Println("enter interceptor 3")
				defer log.Println("exit interceptor 3")
				return do(cli, req)
			}))
	if err != nil {
		fmt.Printf("\nError: %v", err)
		return
	}
	fmt.Printf("\nResponse Status Code: %v", response.RawResponse().StatusCode)
	fmt.Printf("\nResponse Status: %v", response.RawResponse().Status)
	fmt.Printf("\nResponse Header: %v", response.RawResponse().Header)
	bytes, _ := ioutil.ReadAll(response.RawResponse().Body)
	fmt.Printf("\nResponse Body: %s", bytes)
}
