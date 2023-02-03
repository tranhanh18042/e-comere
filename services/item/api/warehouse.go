package api

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type WarehouseRequest struct {
	WarehouseName string `json:"warehouse_name"`
	Address       string `json:"address"`
	PhoneNumber   string `json:"phone_number"`
}

func GetWarehouseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouse model.Warehouse
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-warehouse-by-id",
		}
		err := helper.DBGetWithMetrics(labels, itemDB, &warehouse,
			"SELECT id, warehouse_name, address,phone_number FROM warehouse WHERE id=?",
			ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Info(ctx, "internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: warehouse})
	}
}

func GetListWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouses []model.Warehouse
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-list-warehouse",
		}
		err := helper.DBSelectWithMetrics(labels, itemDB, &warehouses,
			"SELECT * FROM warehouse")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Error(ctx, "item svc get-list-warehouse internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: warehouses})
	}
}

func UpdateWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouseReq WarehouseRequest
		if err := ctx.ShouldBindJSON(&warehouseReq); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			logger.Info(ctx, "invalid params", err)
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "update-warehouse",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB,
			"UPDATE warehouse SET warehouse_name=?, address=?, phone_number=? WHERE id=?",
			warehouseReq.WarehouseName,
			warehouseReq.Address,
			warehouseReq.PhoneNumber,
			ctx.Param("id"),
		)
		if err != nil {
			logger.Info(ctx, "internal error", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: warehouseReq})
	}
}

func CreatWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var warehouseReq WarehouseRequest
		if err := ctx.ShouldBindJSON(&warehouseReq); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			logger.Info(ctx, "invalid params", err)
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "create-warehouse",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB,
			"INSERT INTO warehouse(warehouse_name, address, phone_number) VALUES(?,?,?)",
			warehouseReq.WarehouseName,
			warehouseReq.Address,
			warehouseReq.PhoneNumber)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: warehouseReq})
	}
}
