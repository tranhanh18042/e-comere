package model

import "time"

type Order struct {
	Id             int    `json:"id" db:"id"`
	Status         int    `json:"status" db:"status"`
	CustomerId     int    `json:"customer_id" db:"customer_id"`
	ItemId         int    `json:"item_id" db:"item_id"`
	ItemQuantity   int    `json:"item_quantity" db:"item_quantity"`
	Address        string `json:"address" db:"address"`
	ItemAmount     int    `json:"item_amount" db:"item_amount"`
	ShipFee        int    `json:"ship_fee" db:"ship_fee"`
	TotalAmount    int    `json:"total_amount" db:"total_amount"`
	DiscountAmount int    `json:"discount_amount" db:"discount_amount"`
}
type OrderLog struct {
	OrderId   int       `json:"order_id" db:"order_id"`
	Status    int       `json:"status" db:"status"`
	EventTime time.Time `json:"event_time" db:"event_time"`
}
