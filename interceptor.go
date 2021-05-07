package easyhttp

type Doer func(cli *Client, req *Request) (reply *Reply, err error)

type Interceptor func(cli *Client, req *Request, do Doer) (reply *Reply, err error)

func chainInterceptors(interceptors ...Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	} else if len(interceptors) == 1 {
		return interceptors[0]
	} else {
		return func(httpclient *Client, req *Request, do Doer) (reply *Reply, err error) {
			return interceptors[0](httpclient, req, getChainUnaryInvoker(interceptors, 0, do))
		}
	}
}

func getChainUnaryInvoker(interceptors []Interceptor, curr int, finalInvoker Doer) Doer {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(httpclient *Client, req *Request) (reply *Reply, err error) {
		return interceptors[curr+1](httpclient, req, getChainUnaryInvoker(interceptors, curr+1, finalInvoker))
	}
}
