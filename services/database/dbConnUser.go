package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBUserConn() (dbUser *sql.DB) {
	dbUser, errUser := sql.Open("mysql", "root:root@tcp(db_ecom_user:3307)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
	if errUser != nil {
		panic(errUser)
	}
	return dbUser
}
