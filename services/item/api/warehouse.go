package api

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type Warehouse struct {
	WarehouseName string `json:"WarehouseName"`
	Address       string `json:"Address"`
	PhoneNumber   string `json:"PhoneNumber"`
}
type WarehouseID struct {
	Id            int    `json:"Id"`
	WarehouseName string `json:"WarehouseName"`
	Address       string `json:"Address"`
	PhoneNumber   string `json:"PhoneNumber"`
}

func CreatWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouse Warehouse
		if err := ctx.ShouldBindJSON(&warehouse); err == nil {
			_, err := itemDB.Exec("INSERT INTO warehouse(warehouse_name, address, phone_number) VALUES(?,?,?)", warehouse.WarehouseName, warehouse.Address, warehouse.PhoneNumber)
			if err != nil {
				ctx.JSON(500, gin.H{
					"messages": err,
				})
			}

			ctx.JSON(200, warehouse)
		} else {
			ctx.JSON(500, gin.H{"error": err.Error()})
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
		}
	}
}

func GetWarehouseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		row := itemDB.QueryRow("SELECT id, name_warehouse, address,phone_number FROM warehouse WHERE id= " + ctx.Param("id"))
		var WarehouseID WarehouseID
		if err := row.Scan(&WarehouseID.Id, &WarehouseID.WarehouseName, &WarehouseID.Address, WarehouseID.PhoneNumber); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"messages": "error ",
			})
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricQueryError,
				"env":  "local",
			}).Inc()
			return
		}
		ctx.JSON(200, WarehouseID)
	}
}

func GetWarehouseAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := itemDB.Query("SELECT * FROM warehouse")
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "error",
			})
		}
		var warehouseIDs []WarehouseID

		for rows.Next() {
			var SingleWarehouseID WarehouseID
			if err := rows.Scan(&SingleWarehouseID.Id, &SingleWarehouseID.WarehouseName, &SingleWarehouseID.PhoneNumber, SingleWarehouseID.Address); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"messages": "error",
				})
				metrics.API.ErrCnt.With(prometheus.Labels{
					"svc":  "item",
					"path": ctx.FullPath(),
					"type": helper.MetricQueryError,
					"env":  "local",
				}).Inc()
				return
			}
			warehouseIDs = append(warehouseIDs, SingleWarehouseID)
		}
		ctx.JSON(200, warehouseIDs)
	}
}

func UpdateWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Warehouse Warehouse
		if err := ctx.ShouldBindJSON(&Warehouse); err == nil {
			update, err := itemDB.Prepare("UPDATE warehouse SET warehouse_name=?, address=?, phone_number=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err.Error())
			}
			update.Exec(Warehouse.WarehouseName, Warehouse.Address, Warehouse.PhoneNumber)

			ctx.JSON(200, Warehouse)
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
