# Easyhttp

Easyhttp is a http client written in Golang. The main logic is very simple. Extend functionality using the Chain of Responsibility pattern(inspired by [grpc](https://github.com/grpc/grpc))

## Features

- Easy use.
- Easy to extend via interceptor.
- Built on standard golang `net/http` package.
- Easy intercept and modify HTTP request on-the-fly.
- URL template path params.
- Dependency free.

## Installation

```bash
go get -u github.com/soyacen/easyhttp
```

## Plugins
- [accept](https://github.com/soyacen/easyhttp/tree/main/interceptor/accept) 
- [body](https://github.com/soyacen/easyhttp/tree/main/interceptor/body)
- [auth](https://github.com/soyacen/easyhttp/tree/main/interceptor/auth)
- [client](https://github.com/soyacen/easyhttp/tree/main/interceptor/client)
- [compression](https://github.com/soyacen/easyhttp/tree/main/interceptor/compression)
- [cookie](https://github.com/soyacen/easyhttp/tree/main/interceptor/cookie)
- [download](https://github.com/soyacen/easyhttp/tree/main/interceptor/download)
- [header](https://github.com/soyacen/easyhttp/tree/main/interceptor/header)
- [hystrix](https://github.com/soyacen/easyhttp/tree/main/interceptor/hystrix)
- [logging](https://github.com/soyacen/easyhttp/tree/main/interceptor/logging)
- [multipart](https://github.com/soyacen/easyhttp/tree/main/interceptor/multipart)
- [opentracing](https://github.com/soyacen/easyhttp/tree/main/interceptor/opentracing) 
- [proxy](https://github.com/soyacen/easyhttp/tree/main/interceptor/proxy) 
- [query](https://github.com/soyacen/easyhttp/tree/main/interceptor/query) 
- [retry](https://github.com/soyacen/easyhttp/tree/main/interceptor/retry) 
- [tls](https://github.com/soyacen/easyhttp/tree/main/interceptor/tls) 
- [url](https://github.com/soyacen/easyhttp/tree/main/interceptor/url) 


## Examples

See [examples](https://github.com/soyacen/easyhttp/blob/main/example/get/main.go)


#### Simple request

```go
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

```

#### Send JSON body

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/soyacen/easyhttp"

	easyhttpbody "github.com/soyacen/easyhttp/interceptor/body"
)

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
```


#### Receive JSON Result

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/soyacen/easyhttp"

	easyhttpaccept "github.com/soyacen/easyhttp/interceptor/accpet"
)

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

func main() {
	client := easyhttp.NewClient()
	data := Data{}
	reply, err := client.Get(
		context.Background(),
		"http://httpbin.org/json",
		easyhttp.ChainInterceptor(easyhttpaccept.JSON(&data)))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply.RawResponse().StatusCode)
	fmt.Println(data)
}
```