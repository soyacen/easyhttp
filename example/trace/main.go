package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/soyacen/easyhttp"
)

func main() {
	client := easyhttp.NewClient()
	response, err := client.Trace(
		context.Background(),
		"http://localhost:8080/health")
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
