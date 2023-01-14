package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type ItemDTI struct {
	Quantity    int    `json:"Quantity"`
	Status      int    `json:"Status"`
	ItemName    string `json:"ItemName"`
	UnitPrice   int    `json:"UnitPrice"`
	Description string `json:"description"`
}
type Item struct {
	Id          int    `json:"Id"`
	WarehouseID int    `json:"WarehouseID"`
	ProviderID  int    `json:"ProviderID`
	Quantity    int    `json:"Quantity"`
	Status      int    `json:"Status"`
	ItemName    string `json:"ItemName"`
	UnitPrice   int    `json:"UnitPrice"`
	Description string `json:"description"`
}

func AddItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemDTI ItemDTI
		if err := ctx.ShouldBindJSON(&itemDTI); err == nil {
			_, err := itemDB.Exec("insert into item (warehouse_id,provider_id,quantity,status,item_name, unit_price, description) value(?,?,?,?,?,?,?)", ctx.Param("warehouse_id"), ctx.Param("provider_id"), itemDTI.Quantity, itemDTI.Status, itemDTI.ItemName, itemDTI.UnitPrice, itemDTI.Description)
			if err != nil {
				ctx.JSON(500, gin.H{
					"messages": err,
				})
			}
			ctx.JSON(200, gin.H{
				"messages": "OK",
			})
		} else {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			ctx.JSON(500, gin.H{
				"messages": err.Error(),
			})
		}

	}
}
func GetAllItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := itemDB.Query("select * from item")
		if err != nil {
			panic(err)
		}
		var items []Item
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.Id, &item.WarehouseID, &item.ProviderID, &item.Quantity, &item.Status, &item.ItemName, &item.UnitPrice, item.Description); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err,
				})
				metrics.API.ErrCnt.With(prometheus.Labels{
					"svc":  "item",
					"path": ctx.FullPath(),
					"type": helper.MetricInvalidParams,
					"env":  "local",
				}).Inc()
				return
			}
			items = append(items, item)
		}
		ctx.JSON(200, items)
	}
}
func GetItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var item Item
		row := itemDB.QueryRow("select * from item where id = " + ctx.Param(("id")))
		if err := row.Scan(&item.Id, &item.WarehouseID, &item.ProviderID, &item.Quantity, &item.Status, &item.ItemName, &item.UnitPrice, item.Description); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err,
			})
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			return
		}
		ctx.JSON(200, item)
	}
}
func UpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var item Item
		if err := ctx.ShouldBindJSON(&item); err == nil {
			update, err := itemDB.Preparex("update item set warehouse_id=?,provider_id=?,status=?,item_name=?,unit_price=?,description=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			_, err = update.Exec(item.WarehouseID, item.ProviderID, item.Status, item.ItemName, item.UnitPrice, item.Description)
			if err != nil {
				metrics.DB.ErrCnt.With(prometheus.Labels{
					"env":    "local",
					"type":   "",
					"db":     "ecom_item",
					"target": "error",
				})
			}
			ctx.JSON(200, item)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
		}
	}
}
