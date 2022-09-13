package main

import (
	"database/sql"
	"fmt"
	"golanglearning/new_project/for-gong-zhong-hao/Practical-practice/DBpool"
	"net/http"
)

// var DSN = "username:password@tcp(127.0.0.1:3306)/testdb?charset=utf8"
var DSN = "user=root password=123456 dbname=testdb sslmode=disable"
var poolObj *DBpool.Pool

func main() {
	poolObj = DBpool.NewPool(4, DSN)
	http.HandleFunc("/", test)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err.Error())
	}
}

type user struct {
	Id   int
	Name string
}

func test(writer http.ResponseWriter, req *http.Request) {
	var s user
	res := poolObj.Exec(func(db *sql.DB) interface{} {
		rows, _ := db.Query("select * from \"user\"")
		return rows
	})
	rows := res.(*sql.Rows)
	for rows.Next() {
		err := rows.Scan(&s.Id, &s.Name)
		if err != nil {
			panic(err.Error())
		}
		writer.Write([]byte(fmt.Sprintf("id: %d, name: %s\n", s.Id, s.Name)))
	}
	rows.Close()
}
