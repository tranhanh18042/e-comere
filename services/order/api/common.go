package api

import "github.com/jmoiron/sqlx"

var orderDB *sqlx.DB

// UseDB sets db to be used for service order as global var
func UseDB(db *sqlx.DB) {
	orderDB = db
}
