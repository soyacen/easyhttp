package easyhttp

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

type clientOptions struct {
	interceptors []Interceptor

	//  http.Client field
	transport     http.RoundTripper
	checkRedirect func(req *http.Request, via []*http.Request) error
	jar           http.CookieJar
	timeout       time.Duration
}

func defaultClientOptions() *clientOptions {
	var defaultCookieJar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	o := &clientOptions{
		interceptors:  make([]Interceptor, 0),
		transport:     http.DefaultTransport,
		checkRedirect: nil,
		jar:           defaultCookieJar,
		timeout:       0,
	}
	return o
}

func (o *clientOptions) apply(opts ...ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type ClientOption func(o *clientOptions)

func WithChainInterceptor(interceptors ...Interceptor) ClientOption {
	return func(o *clientOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = transport
	}
}

func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(o *clientOptions) {
		o.jar = jar
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

func WithCheckRedirect(policies ...func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(o *clientOptions) {
		o.checkRedirect = func(req *http.Request, via []*http.Request) error {
			for _, p := range policies {
				if err := p(req, via); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

type Client struct {
	opts      *clientOptions
	rawClient *http.Client
}

func (cli *Client) SetRawClient(rawClient *http.Client) {
	cli.rawClient = rawClient
}

func (cli *Client) RawClient() *http.Client {
	return cli.rawClient
}

func NewClient(opts ...ClientOption) *Client {
	options := defaultClientOptions()
	options.apply(opts...)

	return &Client{
		opts: options,
		rawClient: &http.Client{
			Transport:     options.transport,
			CheckRedirect: options.checkRedirect,
			Jar:           options.jar,
			Timeout:       options.timeout,
		}}
}

func (cli *Client) Get(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodGet, url, opts...)
}

func (cli *Client) Head(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodHead, url, opts...)
}

func (cli *Client) Post(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPost, url, opts...)
}

func (cli *Client) Put(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPut, url, opts...)
}

func (cli *Client) Patch(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPatch, url, opts...)
}

func (cli *Client) Delete(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodDelete, url, opts...)
}

func (cli *Client) Connect(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodConnect, url, opts...)
}

func (cli *Client) Options(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodOptions, url, opts...)
}

func (cli *Client) Trace(ctx context.Context, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodTrace, url, opts...)
}

func (cli *Client) Execute(ctx context.Context, method string, url string, opts ...ExecuteOption) (reply *Reply, err error) {
	options := defaultExecuteOptions()
	options.apply(opts...)
	request := &Request{opts: options}
	var rawReq *http.Request
	if isBodySupported(method) {
		rawReq, err = http.NewRequestWithContext(ctx, method, url, options.body)
		if err != nil {
			return nil, err
		}
	} else {
		rawReq, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	request.rawRequest = rawReq

	interceptors := make([]Interceptor, 0, len(cli.opts.interceptors)+len(request.opts.interceptors))
	for _, interceptor := range cli.opts.interceptors {
		interceptors = append(interceptors, interceptor)
	}
	for _, interceptor := range request.opts.interceptors {
		interceptors = append(interceptors, interceptor)
	}
	request.opts.interceptor = chainInterceptors(interceptors...)

	return cli.execute(request)
}

func (cli *Client) execute(req *Request) (resp *Reply, err error) {
	if req.opts.interceptor != nil {
		return req.opts.interceptor(cli, req, do)
	}
	return do(cli, req)
}

func do(cli *Client, req *Request) (reply *Reply, err error) {
	rawResp, err := cli.rawClient.Do(req.rawRequest)
	if err != nil {
		return nil, err
	}
	reply = &Reply{
		rawRequest:  req.rawRequest,
		rawResponse: rawResp,
	}
	return reply, nil
}
