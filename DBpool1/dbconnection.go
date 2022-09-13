package DBpool1

// https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzkyMzIyNjIxMQ==&action=getalbum&album_id=1795176700377432067&scene=173&from_msgid=2247484559&from_itemidx=1&count=3&nolastread=1#wechat_redirect

/*

 */


type DBConn struct {
	idleTime int  // 标记该数据库连接空闲时间
	Timeout int
}

// 新建数据库连接
func NewDBConn() *DBConn {
	return &DBConn{idleTime: 0}
}

// 关闭数据库连接
func (d *DBConn) Close() {}
