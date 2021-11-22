package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//OpenDB with filepath as connstr
func Open(connstr string) *sql.DB {
	db, err := sql.Open("sqlite3", connstr)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

//CloseDB closes this db
func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println(err)
	}
}
