package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	totalSentRequests := &sync.WaitGroup{}
	allRequestsBackCh := make(chan struct{})
	// 对象
	multiplexCh := make(chan struct {
		result string
		retry  int
	})

	// 启10个goroutine
	for i := 1; i <= 10; i++ {
		// 加一下
		totalSentRequests.Add(1)
		go func() {
			// 标记已经执行完
			defer totalSentRequests.Done()
			// 模拟耗时操作
			time.Sleep(500 * time.Microsecond)
			// 模拟处理成功
			if rand.Intn(500)%2 == 0 {
				multiplexCh <- struct {
					result string
					retry  int
				}{"finsh success", i}
			}
			// 处理失败不关心，当然，也可以加入一个错误的channel中进一步处理
		}()
	}
	for {
		select {
		case <-multiplexCh:
			fmt.Println("finish success")
		case <-allRequestsBackCh:
			// 到这里，说明全部的 goroutine 都执行完毕，但是都请求失败了
			fmt.Println("all req finish，but all fail")
		}

	}

	totalSentRequests.Wait()
	close(allRequestsBackCh)

}
