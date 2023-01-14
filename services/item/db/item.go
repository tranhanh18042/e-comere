package db

import "github.com/jmoiron/sqlx"

func NewItemDB(connStr string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("mysql", connStr)
	return conn, err
}
