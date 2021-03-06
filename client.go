package easyhttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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
	disableCompression  bool
	disableKeepAlives   bool
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
		copy.DisableCompression = o.disableCompression
		copy.DisableKeepAlives = o.disableKeepAlives
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

func WithCertificates(certs ...tls.Certificate) ClientOption {
	return func(o *clientOptions) {
		o.tlsConfig = &tls.Config{Certificates: append([]tls.Certificate{}, certs...)}
	}
}

func WithRootCertificate(pemFilePath string) ClientOption {
	return func(o *clientOptions) {
		rootPemData, err := ioutil.ReadFile(pemFilePath)
		if err != nil {
			panic(err)
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(rootPemData)
		o.tlsConfig = &tls.Config{RootCAs: certPool}
	}
}

func WithRootCertificateFromString(pemContent string) ClientOption {
	return func(o *clientOptions) {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(pemContent))
		o.tlsConfig = &tls.Config{RootCAs: certPool}
	}
}

func WithTLS(config *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConfig = config
	}
}

func WithTLSHandshakeTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.tlsHandshakeTimeout = &timeout
	}
}

func WithProxy(servers map[string]string) ClientOption {
	return func(o *clientOptions) {
		o.proxy = func(req *http.Request) (*url.URL, error) {
			if value, ok := servers[req.URL.Scheme]; ok {
				return url.Parse(value)
			}
			return http.ProxyFromEnvironment(req)
		}
	}
}

func WithDisableCompression() ClientOption {
	return func(o *clientOptions) {
		o.disableCompression = true
	}
}

func WithDisableKeepAlives() ClientOption {
	return func(o *clientOptions) {
		o.disableKeepAlives = true
	}
}

func WithMaxIdleConns(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxIdleConns = &number
	}
}

func WithMaxIdleConnsPerHost(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxIdleConnsPerHost = &number
	}
}

func WithMaxConnsPerHost(number int) ClientOption {
	return func(o *clientOptions) {
		o.maxConnsPerHost = &number
	}
}

func WithIdleConnTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.idleConnTimeout = &timeout
	}
}

func WithWriteBufferSize(size int) ClientOption {
	return func(o *clientOptions) {
		o.writeBufferSize = &size
	}
}

func WithReadBufferSize(size int) ClientOption {
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

func (cli *Client) Get(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodGet, url, itcptrs...)
}

func (cli *Client) Head(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodHead, url, itcptrs...)
}

func (cli *Client) Post(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPost, url, itcptrs...)
}

func (cli *Client) Put(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPut, url, itcptrs...)
}

func (cli *Client) Patch(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodPatch, url, itcptrs...)
}

func (cli *Client) Delete(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodDelete, url, itcptrs...)
}

func (cli *Client) Connect(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodConnect, url, itcptrs...)
}

func (cli *Client) Options(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodOptions, url, itcptrs...)
}

func (cli *Client) Trace(ctx context.Context, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	return cli.Execute(ctx, http.MethodTrace, url, itcptrs...)
}

func (cli *Client) Execute(ctx context.Context, method string, url string, itcptrs ...Interceptor) (reply *Reply, err error) {
	options := defaultExecuteOptions()
	var execOpts []ExecuteOption
	if len(itcptrs) > 0 {
		execOpts = append(execOpts, ChainInterceptor(itcptrs...))
	}
	options.apply(execOpts...)
	request := &Request{opts: options}
	var rawReq *http.Request
	rawReq, err = http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	request.rawRequest = rawReq

	allitcptrs := make([]Interceptor, 0, len(cli.opts.interceptors)+len(request.opts.interceptors))
	for _, itcptr := range cli.opts.interceptors {
		allitcptrs = append(allitcptrs, itcptr)
	}
	for _, interceptor := range request.opts.interceptors {
		allitcptrs = append(allitcptrs, interceptor)
	}
	request.opts.interceptor = chainInterceptors(allitcptrs...)

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
