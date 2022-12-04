package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

func DBItemConn() (dbItem *sqlx.DB) {
	dbItem, errItem := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3309)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
	if errItem != nil {
		panic(errItem)
	}
	return dbItem
}
