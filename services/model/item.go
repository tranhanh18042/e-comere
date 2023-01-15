package model

type Item struct {
	Id          int    `json:"id" db:"id"`
	WarehouseID int    `json:"warehouse_id" db:"warehouse_id"`
	ProviderID  int    `json:"provider_id" db:"provider_id"`
	Quantity    int    `json:"quantity" db:"quantity"`
	Status      int    `json:"status" db:"status"`
	ItemName    string `json:"item_name" db:"item_name"`
	UnitPrice   int    `json:"unit_price" db:"unit_price"`
	Description string `json:"description" db:"description"`
}
