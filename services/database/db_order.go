package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBOrderConn() {
	dbOrder, errOrder := sqlx.Connnect("mysql", "root:root@tcp(db_ecom_order:3308)/ecom_order?collation=utf8mb4_unicode_ci&parseTime=true")
	if errOrder != nil {
		panic(errOrder)
	}
	return dbOrder
}
