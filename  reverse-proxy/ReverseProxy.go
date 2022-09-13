package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)


/*
	反向代理使用场景
	1.负载均衡（Load balancing）：反向代理可以提供负载均衡解决方案，将传入的流量均匀地分布在不同的服务器之间，以防止单个服务器过载。
	2.防止安全攻击：由于真正的后端服务器永远不需要暴露公共 IP，所以 DDoS 等攻击只能针对反向代理进行， 这能确保在网络攻击中尽量多的保护你的资源，真正的后端服务器始终是安全的。
	3.缓存：假设你的实际服务器与用户所在的地区距离比较远，那么你可以在当地部署反向代理，它可以缓存网站内容并为当地用户提供服务。
	4.SSL 加密：由于与每个客户端的 SSL 通信会耗费大量的计算资源，因此可以使用反向代理处理所有与 SSL 相关的内容，然后释放你真正服务器上的宝贵资源。
 */

// NewProxy方法构建input targetHost对象，构建反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	// 修改请求
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	// 修改回应
	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()


	return proxy, nil
}

// ProxyRequestHandler方法 则是使用proxy处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}

}

// HttpUtil 反向代理为我们提供了一种非常简单的机制来修改我们从服务器获得的响应， 可以根据你的应用场景来缓存或更改此响应
// 设置了自定义 Header头。同样也可以读取响应体正文，并对其进行更改或缓存，然后将其设置回客户端。
func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

// 在 modifyResponse 中，可以返回一个错误（如果你在处理响应发生了错误），
// 如果你设置了 proxy.ErrorHandler,
// modifyResponse 返回错误时会自动调用 ErrorHandler 进行错误处理。
func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return errors.New("response body is invalid")
	}
}



// 修改请求
// 可以在将请求发送到服务器之前对其进行修改。在下面的例子中，我们将会在请求发送到服务器之前添加了一个 Header 头。
func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}




func main() {
	// 到达我们代理服务器的任何请求都会被代理到位于 http://my-api-server.com
	targetHost := "http://my-api-server.com"
	proxy, err := NewProxy(targetHost)
	if err != nil {
		panic(err)
	}

	// 使用Proxy处理所有请求
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe("8080", nil))



}