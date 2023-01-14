package api

import "github.com/jmoiron/sqlx"

var customerDB *sqlx.DB

// UseDB sets db to be used for service order as global var
func UseDB(db *sqlx.DB) {
	customerDB = db
}
