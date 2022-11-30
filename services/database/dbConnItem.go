package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBItemConn() (dbItem *sql.DB) {
	dbItem, errItem := sql.Open("mysql", "root:root@tcp(db_ecom_item:3309)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
	if errItem != nil {
		panic(errItem)
	}
	return dbItem
}
