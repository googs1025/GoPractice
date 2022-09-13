package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 服务端收到请求之后就会从这个Reader中调用Read()函数去读取数据，通常情况当服务端去读取数据的时候，
// offset会随之改变，下一次再读的时候会从offset位置继续向后读取。所以如果直接重试，会出现读不到 Reader的情况。

func main() {

	// 启一个http监听服务
	go func() {
		//
		http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 停止20秒
			time.Sleep(time.Second * 3)
			// 读 body
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Printf("received body with length %v containing: %v\n", len(body), string(body))
			// 写回http响应码
			w.WriteHeader(http.StatusOK)
		}))
		// 监听
		http.ListenAndServe(":8090", nil)
	}()
	fmt.Print("Try with bare strings.Reader\n")

	// 要发送的body
	originalBody := []byte("abcdefghigklmnopqrst")
	// 创一个新的读取器
	reader := strings.NewReader(string(originalBody))
	req, _ := http.NewRequest("POST", "http://localhost:8090/", reader)

	// 调用retryDo
	retryDo(req, 10)
}

func retryDo(req *http.Request, maxRetries int) (*http.Response, error) {

	// 建立客户端，设置超时时间
	client := http.Client{
		Timeout: time.Second * 2,
	}

	// 重试次数
	AttemptLimit := maxRetries
	if AttemptLimit <= 0 {
		AttemptLimit = 1
	}

	// 变量
	var (
		originalBody []byte
		err          error
	)

	// 如果请求不为空，需要先复制一下！
	if req != nil && req.Body != nil {
		originalBody, err = copyBody(req.Body)
		resetBody(req, originalBody)
	}

	// 请求
	var resp *http.Response

	for i := 0; i < AttemptLimit; i++ {
		// 调用客户端发送，并接收返回的请求
		resp, err = client.Do(req)
		// 错误处理
		if err != nil {
			fmt.Printf("error sending the first time: %v\n", err)
		}
		// 重试 500 以上的错误码
		if err == nil && resp.StatusCode < 500 {
			return resp, err
		}
		// 如果正在重试，那么释放fd
		if resp != nil {
			resp.Body.Close()
		}
		// 重置body
		if req.Body != nil {
			resetBody(req, originalBody)
		}

	}
	// 如果到这里就代表重试也不行
	return resp, req.Context().Err()
}


func copyBody(src io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	src.Close()
	return b, nil
}

// 当再次请求的时候，发现 client 请求的 Body 数据并不是我们预期的20个长度，
// 而是 0，导致了 err。因此需要将Body这个Reader 进行重置
func resetBody(request *http.Request, originalBody []byte) {
	request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
}



