package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBUserConn() {
	dbUser, errUser := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3307)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
	if errUser != nil {
		panic(errUser)
	}
	return dbUser
}
