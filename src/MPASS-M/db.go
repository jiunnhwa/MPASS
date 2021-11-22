package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB //global singleton for connection pool reuse

const Connstr = "user:password@tcp(192.168.0.162:3306)/pegasus?parseTime=true"

func OpenDB(connstr string) {
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
	}
	DB = db
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		panic(err.Error())
	}
}
