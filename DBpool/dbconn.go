package DBpool

import (
	"database/sql"
	"fmt"
	"time"
)

// Golang 实现简单的数据库连接池
// https://www.bilibili.com/read/cv14785530/


type DBConn struct {
	ID	int
	IdleTime *time.Timer
	Timeout int
	Conn 	*sql.DB

}

// 执行的sql语句
type ExecFunc func(db *sql.DB) interface{}

// 执行时，需要计时
func (dbconn *DBConn) exec(fn ExecFunc) interface{} {
	dbconn.countDown()
	return fn(dbconn.Conn)
}

// 定时关闭conn
func (dbconn *DBConn) countDown() {
	if dbconn.IdleTime != nil {
		dbconn.IdleTime.Stop()
	}
	dbconn.IdleTime = time.AfterFunc(time.Duration(dbconn.Timeout*int(time.Second)), func() {
		dbconn.Conn.Close()
		dbconn.Conn = nil
		fmt.Printf("Connection %d 关闭\n", dbconn.ID)
	})
}

// 建立conn
func (dbconn *DBConn) open(driver string, dsn string) {

	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err.Error())
	}
	dbconn.Conn = db

}

// 构建函数
func NewDBConn(id int, driver string, dsn string) *DBConn {
	conn := &DBConn{
		ID: id,
		Timeout: 30,
	}
	conn.open(driver, dsn)
	return conn

}


