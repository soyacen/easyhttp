package main

import (
	"context"
	"fmt"
	"log"

	"github.com/soyacen/easyhttp"

	easyhttpresp "github.com/soyacen/easyhttp/interceptor/resp"
)

func main() {
	client := easyhttp.NewClient()
	data := Data{}
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/json",
		easyhttp.ChainInterceptor(easyhttpresp.JSON(&data)))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply.RawResponse().StatusCode)
	fmt.Println(data)

	fmt.Println("======================")

	data1 := Data{}
	reply1, err := client.Get(
		context.Background(),
		"http://httpbin.org/gzip",
		easyhttp.ChainInterceptor(easyhttpresp.JSON(&data1)))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply1.RawResponse().StatusCode)
	fmt.Println(data1)

	fmt.Println("======================")

	data2 := Data{}
	reply2, err := client.Get(
		context.Background(),
		"http://httpbin.org/xml",
		easyhttp.ChainInterceptor(easyhttpresp.XML(&data2)))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply2.RawResponse().StatusCode)
	fmt.Println(data2)
}

type Data struct {
	Slideshow Slideshow `json:"slideshow" xml:"slideshow"`
}

type Slideshow struct {
	Author string  `json:"author" xml:"author"`
	Date   string  `json:"date" xml:"date"`
	Slides []Slide `json:"slides" xml:"slides"`
	Title  string  `json:"title" xml:"title"`
}

type Slide struct {
	Title string   `json:"title" xml:"title"`
	Type  string   `json:"type" xml:"type"`
	Items []string `json:"items,omitempty" xml:"items"`
}
