package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/soyacen/easyhttp"
	easyhttpbalancer "github.com/soyacen/easyhttp/interceptor/balancer"
)

func RoundRobin() {
	client := easyhttp.NewClient(
		easyhttp.WithChainInterceptor(
			easyhttpbalancer.Interceptor(
				easyhttpbalancer.WithPicker(easyhttpbalancer.NewRoundRobinPicker()),
			),
		))

	for i := 0; i < 10; i++ {
		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		response, err := client.Get(timeout, "http://httpbin.org:80,google.com:80/get")
		if err != nil {
			fmt.Printf("\nError: %v", err)
			continue
		}
		fmt.Printf("\nResponse Status Code: %v", response.RawResponse().StatusCode)
		fmt.Printf("\nResponse Status: %v", response.RawResponse().Status)
		fmt.Printf("\nResponse Header: %v", response.RawResponse().Header)
		bytes, _ := ioutil.ReadAll(response.RawResponse().Body)
		fmt.Printf("\nResponse Body: %s", bytes)
	}
}

func Rand() {
	client := easyhttp.NewClient(
		easyhttp.WithChainInterceptor(
			easyhttpbalancer.Interceptor(
				easyhttpbalancer.WithPicker(easyhttpbalancer.NewRandomPicker()),
			),
		))
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	response, err := client.Get(timeout, "http://httpbin.org:80,google.com:80/get")
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

func Hash() {
	client := easyhttp.NewClient(
		easyhttp.WithChainInterceptor(
			easyhttpbalancer.Interceptor(
				easyhttpbalancer.WithPicker(easyhttpbalancer.NewHashPicker("token", easyhttpbalancer.Query)),
			),
		))
	response, err := client.Get(context.Background(), "http://httpbin.org:80,google.com:80/get?token=100023456")
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

func First() {
	client := easyhttp.NewClient(
		easyhttp.WithChainInterceptor(
			easyhttpbalancer.Interceptor(
				easyhttpbalancer.WithPicker(easyhttpbalancer.NewFirstPicker()),
			),
		))
	response, err := client.Get(context.Background(), "http://httpbin.org:80,google.com:80/get")
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
