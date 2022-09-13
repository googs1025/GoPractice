package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	go func() {
		http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * 20)
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Printf("received body with length %v containing: %v\n", len(body), string(body))
			w.WriteHeader(http.StatusOK)
		}))
		http.ListenAndServe(":8090", nil)
	}()
	fmt.Print("Try with bare strings.Reader\n")
	retryDo(req)
}


func retryDo(req *http.Request, maxRetries int, timeout time.Duration,
	backoffStrategy BackoffStrategy) (*http.Response, error) {
	var (
		originalBody []byte
		err          error
	)
	if req != nil && req.Body != nil {
		originalBody, err = copyBody(req.Body)
		resetBody(req, originalBody)
	}
	if err != nil {
		return nil, err
	}
	AttemptLimit := maxRetries
	if AttemptLimit <= 0 {
		AttemptLimit = 1
	}

	client := http.Client{
		Timeout: timeout,
	}
	var resp *http.Response
	//重试次数
	for i := 1; i <= AttemptLimit; i++ {
		resp, err = client.Do(req)
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
		time.Sleep(backoffStrategy(i) + 1*time.Microsecond)
	}
	// 到这里，说明重试也没用
	return resp, req.Context().Err()
}

func copyBody(src io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, ErrReadingRequestBody
	}
	src.Close()
	return b, nil
}

func resetBody(request *http.Request, originalBody []byte) {
	request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
}