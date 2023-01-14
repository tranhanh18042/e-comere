package api

import "github.com/tranhanh18042/e-comere/services/helper"

var customerDB *helper.SvcDB

// UseDB sets db to be used for service order as global var
func UseDB(db *helper.SvcDB) {
	customerDB = db
}
