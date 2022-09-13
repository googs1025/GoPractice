package DBpool1

import "sync"

type Pool struct {
	mu      sync.Mutex
	minConn int // 最小连接数
	maxConn int // 最大连接数
	numConn int // 池已申请的连接数
	conns   chan *DBConn //当前池中空闲连接实例
	close   bool
}

// 初始化池实例
func NewPool(min, max int) *Pool {
	p := &Pool{
		minConn: min,
		maxConn: max,
		numConn: min,
		conns:   make(chan *DBConn, max),
		close:   false,
	}
	for i := 0; i < min; i++ {
		p.conns <- NewDBConn()
	}
	return p
}
