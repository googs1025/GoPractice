package DBpool

type Pool struct {

	numConn int
	conns	chan *DBConn
	Driver string
	Dsn    string

}

func NewPool(size int, dsn string) *Pool {
	// 池对象
	p := &Pool{
		numConn: size,
		conns: make(chan *DBConn, size),
		Driver: "postgres",
		Dsn: dsn,

	}

	// 建立对象放入池中
	for i := 0; i < p.numConn; i++ {
		p.conns <-NewDBConn(i, p.Driver, p.Dsn)
	}

	return p
}


