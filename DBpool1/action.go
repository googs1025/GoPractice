package DBpool1

import "fmt"

// 从池中取出连接
func (p *Pool) Get() *DBConn {
	if p.close {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果连接池中的连接数超过最大连接数且chan中有可用的资源，直接拿取。
	if p.numConn >= p.maxConn || len(p.conns) > 0 {

		d := <-p.conns // 若池中没有可取的连接，则等待其他请求返回连接至池中再取
		return d
	}
	p.numConn++
	return NewDBConn() //申请新的连接
}

// 将连接返回池中
func (p *Pool) Put(d *DBConn) {
	if p.close {
		fmt.Println("连接池已关闭")
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conns <- d
}

// 关闭池
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for d := range p.conns {
		d.Close()
	}
	p.close = true
}
