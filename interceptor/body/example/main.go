package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/soyacen/easyhttp"
	easyhttpbody "github.com/soyacen/easyhttp/interceptor/form"
)

func ExampleJSONString() {
	client := easyhttp.NewClient()
	reply, err := client.Post(
		context.Background(),
		"http://httpbin.org/post",
		easyhttp.ChainInterceptor(easyhttpbody.JSON(`{"extra":null,"accesskey":"rfe65iisgjr4ltgp","expid":"62","entity":"869791045881921","traceid":"f7eef0d861379d6681940f07b548eff1","bucket":"2499","group":"174","ts":"1620389233"}`)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func ExampleXML() {
	client := easyhttp.NewClient()
	type Data struct {
		Comments  string `xml:"comments"`
		Custemail string `xml:"custemail"`
		Custname  string `xml:"custname"`
		Custtel   string `xml:"custtel"`
		Delivery  string `xml:"delivery"`
		Size      string `xml:"size"`
		Topping   string `xml:"topping"`
	}
	data := Data{
		Comments:  "adsa",
		Custemail: "dsadsa@ww",
		Custname:  "sad",
		Custtel:   "dsad",
		Delivery:  "18:30",
		Size:      "small",
		Topping:   "bacon",
	}
	reply, err := client.Post(
		context.Background(),
		"http://httpbin.org/post",
		easyhttp.ChainInterceptor(easyhttpbody.XML(data)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func ExampleJSON() {
	client := easyhttp.NewClient()
	data := map[string]string{
		"comments":  "adsa",
		"custemail": "dsadsa@ww",
		"custname":  "sad",
		"custtel":   "dsad",
		"delivery":  "18:30",
		"size":      "small",
		"topping":   "bacon",
	}
	reply, err := client.Post(
		context.Background(),
		"http://httpbin.org/post",
		easyhttp.ChainInterceptor(easyhttpbody.JSON(data)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))
}

func ExampleForm() {
	client := easyhttp.NewClient()
	form := url.Values{
		"comments":  []string{"adsa"},
		"custemail": []string{"dsadsa@ww"},
		"custname":  []string{"sad"},
		"custtel":   []string{"dsad"},
		"delivery":  []string{"18:30"},
		"size":      []string{"small"},
		"topping":   []string{"bacon"},
	}
	reply, err := client.Post(
		context.Background(),
		"http://httpbin.org/post",
		easyhttp.ChainInterceptor(easyhttpbody.Form(form)))
	if err != nil {
		log.Fatalln(err)
	}
	response, err := httputil.DumpResponse(reply.RawResponse(), true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(response))

}
