package DBpool



func (p *Pool) Exec(queryFunc ExecFunc) interface{} {

	// get一个池中资源
	connection := <-p.conns
	// 如果没有，建立一个
	if connection.Conn == nil {
		connection.open(p.Driver, p.Dsn)
	}
	// 执行命令
	result := connection.exec(queryFunc)

	// 执行后需要放回池中
	defer func(connection *DBConn) {
		p.conns <-connection
	}(connection)

	return result

}




