package api

import "github.com/jmoiron/sqlx"

var itemDB *sqlx.DB

// UseDB sets db to be used for service item as global var
func UseDB(db *sqlx.DB) {
	itemDB = db
}
