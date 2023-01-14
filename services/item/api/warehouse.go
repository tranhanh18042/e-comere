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
	Name_warehouse string `json:"name_warehouse"`
	Address        string `json:"address"`
	Phone_Number   string `json:"phone_number"`
}
type WarehouseId struct {
	Id             int    `json:"id"`
	Name_warehouse string `json:"name_warehouse"`
	Address        string `json:"address"`
	Phone_Number   string `json:"phone_number"`
}

func CreatWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouse Warehouse
		if err := ctx.ShouldBindJSON(&warehouse); err == nil {
			_, err := itemDB.Exec("INSERT INTO warehouse(name_warehouse, address, phone_number) VALUES(?,?,?)", warehouse.Name_warehouse, warehouse.Address, warehouse.Phone_Number)
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
		var warehouseId WarehouseId
		if err := row.Scan(&warehouseId.Id, &warehouseId.Name_warehouse, &warehouseId.Address, warehouseId.Phone_Number); err == nil {
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
		ctx.JSON(200, warehouseId)
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
		var warehouseId []WarehouseId

		for rows.Next() {
			var singleWarehouse WarehouseId
			if err := rows.Scan(&singleWarehouse.Id, &singleWarehouse.Name_warehouse, &singleWarehouse.Phone_Number, singleWarehouse.Address); err == nil {
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
			warehouseId = append(warehouseId, singleWarehouse)
		}
		ctx.JSON(200, warehouseId)
	}
}

func UpdateWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouse Warehouse
		if err := ctx.ShouldBindJSON(&warehouse); err == nil {
			update, err := itemDB.Prepare("UPDATE warehouse SET name_warehouse=?, address=?, phone_number=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err.Error())
			}
			update.Exec(warehouse.Name_warehouse, warehouse.Address, warehouse.Phone_Number)

			ctx.JSON(200, warehouse)
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
