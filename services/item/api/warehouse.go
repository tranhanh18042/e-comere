package api

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type WarehouseRequest struct {
	WarehouseName string `json:"warehouse_name"`
	Address       string `json:"address"`
	PhoneNumber   string `json:"phone_number"`
}

func CreatWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouseReq WarehouseRequest
		if err := ctx.ShouldBindJSON(&warehouseReq); err == nil {
			_, err := itemDB.Exec("INSERT INTO warehouse(warehouse_name, address, phone_number) VALUES(?,?,?)",
			warehouseReq.WarehouseName,
			warehouseReq.Address,
			warehouseReq.PhoneNumber)
			if err != nil {
				ctx.JSON(500, gin.H{
					"messages": err,
				})
			}

			ctx.JSON(200, warehouseReq)
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
		var warehouse model.Warehouse
		if err := row.Scan(&warehouse.Id, &warehouse.WarehouseName, &warehouse.Address, warehouse.PhoneNumber); err == nil {
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
		ctx.JSON(200, warehouse)
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
		var warehouses []model.Warehouse
		for rows.Next() {
			var singleWarehouseID model.Warehouse
			if err := rows.Scan(&singleWarehouseID.Id, &singleWarehouseID.WarehouseName, &singleWarehouseID.PhoneNumber, singleWarehouseID.Address); err == nil {
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
			warehouses = append(warehouses, singleWarehouseID)
		}
		ctx.JSON(200, warehouses)
	}
}

func UpdateWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouseReq WarehouseRequest
		if err := ctx.ShouldBindJSON(&warehouseReq); err == nil {
			update, err := itemDB.Prepare("UPDATE warehouse SET warehouse_name=?, address=?, phone_number=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err.Error())
			}
			update.Exec(warehouseReq.WarehouseName, warehouseReq.Address, warehouseReq.PhoneNumber)

			ctx.JSON(200, warehouseReq)
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
