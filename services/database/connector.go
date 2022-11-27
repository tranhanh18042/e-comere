package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBconnection() (db *sqlx.DB) {
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3307)/usera")
	if err != nil {
		panic(err)
	}
	return db
}
