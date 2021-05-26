package easyhttprespbody

import (
	"context"
	"log"
	"testing"

	"github.com/soyacen/easyhttp"
)

func TestRespBody(t *testing.T) {
	client := easyhttp.NewClient()
	data := Data{}
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/json",
		easyhttp.ChainInterceptor(JSON(&data)))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply.RawResponse().StatusCode)
	t.Log(data)
	t.Log("======================")

	data1 := Data{}
	reply1, err := client.Get(
		context.Background(),
		"http://httpbin.org/gzip",
		easyhttp.ChainInterceptor(JSON(&data1)))
	if err != nil {
		log.Fatalln(err)
	}
	t.Log(reply1.RawResponse().StatusCode)
	t.Log(data1)

	t.Log("======================")

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
