package api

import "github.com/tranhanh18042/e-comere/services/helper"

var itemDB *helper.SvcDB

// UseDB sets db to be used for service item as global var
func UseDB(db *helper.SvcDB) {
	itemDB = db
}
