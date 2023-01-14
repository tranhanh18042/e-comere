package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
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
	Id          int    `json:"id" db:"id"`
	WarehouseID int    `json:"warehouse_id" db:"warehouse_id"`
	ProviderID  int    `json:"provider_id" db:"provider_id"`
	Quantity    int    `json:"quantity" db:"quantity"`
	Status      int    `json:"status" db:"status"`
	ItemName    string `json:"item_name" db:"item_name"`
	UnitPrice   int    `json:"unit_price" db:"unit_price"`
	Description string `json:"description" db:"description"`
}

func GetAllItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var items []Item
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-all-items",
		}
		err := helper.DBSelectWithMetrics(labels, itemDB.DB, &items, "select * from item")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Error(ctx, "item svc get-all-items internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(200, helper.SuccessResponse{
			Payload: items,
		})
	}
}
func GetItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var item Item
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-item-by-id",
		}
		err := helper.DBGetWithMetrics(labels, itemDB, &item, "select * from item where id = ?", ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(200, helper.SuccessResponse{
			Payload: item,
		})
	}
}

func AddItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemDTI ItemDTI
		if err := ctx.ShouldBindJSON(&itemDTI); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}

		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "add-item",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB,
			"INSERT INTO item(warehouse_id, provider_id, quantity, status,item_name, unit_price, description) VALUE(?,?,?,?,?,?,?)",
			ctx.Param("warehouse_id"),
			ctx.Param("provider_id"),
			itemDTI.Quantity,
			itemDTI.Status,
			itemDTI.ItemName,
			itemDTI.UnitPrice,
			itemDTI.Description)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: itemDTI})
	}
}

func UpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemUpdate Item
		if err := ctx.ShouldBindJSON(&itemUpdate); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()

			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}

		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "update-item",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB,
			"update item set warehouse_id=?, provider_id=?, status=?, item_name=?, unit_price=?, description=? where id= ?",
			itemUpdate.WarehouseID,
			itemUpdate.ProviderID,
			itemUpdate.Status,
			itemUpdate.ItemName,
			itemUpdate.UnitPrice,
			itemUpdate.Description,
			ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(200, itemUpdate)
	}
}
