package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type ItemDTI struct {
	Amount      int    `json:"Amount"`
	Status      int    `json:"Status"`
	Name_item   string `json:"Name_item"`
	Price       int    `json:"Price"`
	Description string `json:"description"`
}
type Item struct {
	Id          int    `json:"Id"`
	WarehouseId int    `json:"WarehouseId"`
	ProviderId  int    `json:"Provider`
	Amount      int    `json:"Amount"`
	Status      int    `json:"Status"`
	Name_item   string `json:"Name_item"`
	Price       int    `json:"Price"`
	Description string `json:"description"`
}

type ItemDetail struct {
	Id int `json:"Id"`
	WarehouseDTO
	ProviderDTO
	Amount      int    `json:"Amount"`
	Status      int    `json:"Status"`
	Name_item   string `json:"Name_item"`
	Price       int    `json:"Price"`
	Description string `json:"description"`
}

type WarehouseDTO struct {
	Id             int    `json:"Id"`
	Name_warehouse string `json:"Name_warehouse"`
	Address        string `json:"Address"`
	Phone_number   string `json:"Phone_number"`
}
type ProviderDTO struct {
	Id            int    `json:"Id"`
	Name_provider string `json:"Name_warehouse"`
	Address       string `json:"Address"`
	Phone_number  string `json:"Phone_number"`
}

func AddItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemDTI ItemDTI
		if err := ctx.ShouldBindJSON(&itemDTI); err == nil {
			_, err := itemDB.Exec("insert into item (warehouse_id,provider_id,amount,status,name_item, price, description) value(?,?,?,?,?,?,?)", ctx.Param("warehouse_id"), ctx.Param("provider_id"), itemDTI.Amount, itemDTI.Status, itemDTI.Name_item, itemDTI.Price, itemDTI.Description)
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
			if err := rows.Scan(&item.Id, &item.WarehouseId, &item.ProviderId, &item.Amount, &item.Status, &item.Name_item, &item.Price, item.Description); err == nil {
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
		if err := row.Scan(&item.Id, &item.WarehouseId, &item.ProviderId, &item.Amount, &item.Status, &item.Name_item, &item.Price, item.Description); err == nil {
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
			update, err := itemDB.Preparex("update item set warehouse_id=?,provider_id=?,status=?,name_item=?,price=?,description=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			_, err = update.Exec(item.WarehouseId, item.ProviderId, item.Status, item.Name_item, item.Price, item.Description)
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
