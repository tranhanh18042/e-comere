package model

type Warehouse struct {
	Id            int    `json:"id" db:"id"`
	WarehouseName string `json:"warehouse_name" db:"warehouse_name"`
	Address       string `json:"address" db:"address"`
	PhoneNumber   string `json:"phone_number" db:"phone_number"`
}
