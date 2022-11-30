package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBOrderConn() (dbOrder *sql.DB) {
	dbOrder, errOrder := sql.Open("mysql", "root:root@tcp(db_ecom_order:3308)/ecom_order?collation=utf8mb4_unicode_ci&parseTime=true")
	if errOrder != nil {
		panic(errOrder)
	}
	return dbOrder
}
