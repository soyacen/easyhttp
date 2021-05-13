package easyhttp

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

type clientOptions struct {
	interceptors []Interceptor

	//  http.Client field
	transport           http.RoundTripper
	checkRedirect       func(req *http.Request, via []*http.Request) error
	jar                 http.CookieJar
	timeout             *time.Duration
	tlsConfig           *tls.Config
	tlsHandshakeTimeout *time.Duration
	proxy               func(req *http.Request) (*url.URL, error)
	compression         *bool
	keepAlives          *bool
	maxIdleConns        *int
	maxIdleConnsPerHost *int
	maxConnsPerHost     *int
	idleConnTimeout     *time.Duration
	writeBufferSize     *int
	readBufferSize      *int
}

func defaultClientOptions() *clientOptions {
	var defaultCookieJar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	o := &clientOptions{
		interceptors:  make([]Interceptor, 0),
		transport:     nil,
		checkRedirect: nil,
		jar:           defaultCookieJar,
		timeout:       new(time.Duration),
	}
	return o
}

func (o *clientOptions) apply(opts ...ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
	if o.transport == nil {
		transport := http.DefaultTransport.(*http.Transport)
		copy := *transport
		if o.tlsConfig != nil {
			copy.TLSClientConfig = o.tlsConfig
		}
		if o.tlsHandshakeTimeout != nil {
			copy.TLSHandshakeTimeout = *o.tlsHandshakeTimeout
		}
		if o.proxy != nil {
			copy.Proxy = o.proxy
		}
		if o.compression != nil {
			copy.DisableCompression = !*o.compression
		}
		if o.keepAlives != nil {
			copy.DisableKeepAlives = !*o.keepAlives
		}
		if o.maxIdleConns != nil {
			copy.MaxIdleConns = *o.maxIdleConns
		}
		if o.maxIdleConnsPerHost != nil {
			copy.MaxIdleConnsPerHost = *o.maxIdleConnsPerHost
		}
		if o.maxConnsPerHost != nil {
			copy.MaxConnsPerHost = *o.maxConnsPerHost
		}
		if o.idleConnTimeout != nil {
			copy.IdleConnTimeout = *o.idleConnTimeout
		}
		if o.writeBufferSize != nil {
			copy.WriteBufferSize = *o.writeBufferSize
		}
		if o.readBufferSize != nil {
			copy.ReadBufferSize = *o.readBufferSize
		}
		o.transport = &copy
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
		o.timeout = &timeout
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

func TLS(config *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConfig = config
	}
}

func TLSHandshakeTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.tlsHandshakeTimeout = &timeout
	}
}

func Proxy(servers map[string]string) ClientOption {
	return func(o *clientOptions) {
		o.proxy = func(req *http.Request) (*url.URL, error) {
			if value, ok := servers[req.URL.Scheme]; ok {
				return url.Parse(value)
			}
			return http.ProxyFromEnvironment(req)
		}
	}
}

func Compression(enabled bool) ClientOption {
	return func(o *clientOptions) {
		o.compression = &enabled
	}
}

func KeepAlives(enabled bool) ClientOption {
	return func(o *clientOptions) {
		o.keepAlives = &enabled
	}
}

func MaxIdleConns(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxIdleConns = &number
	}
}

func MaxIdleConnsPerHost(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxIdleConnsPerHost = &number
	}
}

func MaxConnsPerHost(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxConnsPerHost = &number
	}
}

func IdleConnTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.idleConnTimeout = &timeout
	}
}

func WriteBufferSize(size int) ClientOption {
	return func(o *clientOptions) {
		o.writeBufferSize = &size
	}
}

func ReadBufferSize(size int) ClientOption {
	return func(o *clientOptions) {
		o.readBufferSize = &size
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
			Timeout:       *options.timeout,
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
	rawReq, err = http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
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
