package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/soyacen/easyhttp"

	easyhttpdownload "github.com/soyacen/easyhttp/interceptor/download"
)

func main() {
	filename := "/tmp/httpbin.webp"
	os.Remove(filename)
	client := easyhttp.NewClient()
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/image",
		easyhttp.ChainInterceptor(easyhttpdownload.Download(filename)),
	)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply.RawResponse().StatusCode)

	fileInfo, _ := os.Stat(filename)
	fmt.Println(fileInfo.Size())
}
