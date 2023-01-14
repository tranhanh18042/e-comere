package db

import "github.com/jmoiron/sqlx"

func NewCustomerDB(connStr string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("mysql", connStr)
	return conn, err
}
