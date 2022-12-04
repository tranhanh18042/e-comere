package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

func DBOrderConn() (dbOrder *sqlx.DB) {
	dbOrder, errOrder := sqlx.Connect("mysql", "root:root@tcp(db_ecom_order:3308)/ecom_order?collation=utf8mb4_unicode_ci&parseTime=true")
	if errOrder != nil {
		panic(errOrder)
	}
	return dbOrder
}
