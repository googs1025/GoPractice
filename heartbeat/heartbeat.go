package main

import (
	"context"
	"time"
	"io"
	"fmt"
)

/*
    https://mp.weixin.qq.com/s/3Q6P8ikf5gePbeUGwH4PBQ
	当节点间有长时间连接时，需要确认节点间的存活状态，最好的方式就是心跳机制。
    一个心跳消息需要被发送到远端服务并等待回复，我们可以根据回复情况提前终止连接。节点将以一定间隔时间来发送消息，类似心跳。

    实现方式：
	需要一个goroutine定期到发送ping消息。
	如果最近收到远程服务的回复，就不需要发送不必要的ping消息，因此需要可以重置ping计时器功能。
 */

// 默认间隔

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration, defaultPingInterval time.Duration) {
	// 间隔
	var interval time.Duration
	// 用select通道
	select {
	case <-ctx.Done():	//
		return
	case interval = <-reset: //读取更新的心跳间隔时间
	default:
	}
	//  设置时间
	if interval < 0 {
		interval = defaultPingInterval
	}
	// 创一个计时器
	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("pingddddd")); err != nil {
				//在此跟踪并执行连续超时
				return
			}
		}
		_ = timer.Reset(interval) //重制心跳上报时间间隔
	}
}

func ExamplePinger(defaultPingInterval time.Duration) {
	// 取消控制
	ctx, cancelFunc := context.WithCancel(context.Background())
	// ？？？？
	r, w := io.Pipe() //代替网络连接net.Conn
	// 通知结束chan
	done := make(chan struct{})
	// 放时间的chan
	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second //ping间隔初始值

	// 启一个goroutine 调用心跳功能
	go func() {
		Pinger(ctx, w, resetTimer, defaultPingInterval)
		close(done)	// 关闭通道
	}()
	// 接收心跳消息
	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("resetting time (%s)\n", d)
			resetTimer <- d
		}

		now := time.Now()
		// 读取的方式：订制一块缓冲区，用它来读取
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))
	}
	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		fmt.Printf("Run %d\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}
	// 取消控制
	cancelFunc() //取消context使pinger退出
	<-done	// 确认goroutine关闭
}

func main()  {
	// 调用
	ExamplePinger(time.Second * 3)

}