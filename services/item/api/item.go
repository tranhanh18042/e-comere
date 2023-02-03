package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type ItemRequest struct {
	Quantity    int    `json:"quantity"`
	Status      int    `json:"status"`
	ItemName    string `json:"item_name"`
	UnitPrice   int    `json:"unit_price"`
	Description string `json:"description"`
	WarehouseID int    `json:"warehouse_id"`
	ProviderId  int    `json:"provider_id"`
}

func GetListItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var items []model.Item
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-list-items",
		}
		err := helper.DBSelectWithMetrics(labels, itemDB, &items, "SELECT * FROM item")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Error(ctx, "item svc get-list-items internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{
			Payload: items,
		})
	}
}
func GetItemByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var item model.Item
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-item-by-id",
		}
		err := helper.DBGetWithMetrics(labels, itemDB, &item, "SELECT * FROM item WHERE id = ?", ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{
			Payload: item,
		})
	}
}

func CreateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemReq ItemRequest
		if err := ctx.ShouldBindJSON(&itemReq); err != nil {
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
			Target: "create-item",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB,
			"INSERT INTO item(warehouse_id, provider_id, quantity, status,item_name, unit_price, description) VALUE(?,?,?,?,?,?,?)",
			itemReq.WarehouseID,
			itemReq.ProviderId,
			itemReq.Quantity,
			itemReq.Status,
			itemReq.ItemName,
			itemReq.UnitPrice,
			itemReq.Description)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: itemReq})
	}
}

func UpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemUpdate model.Item
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
			"UPDATE item SET warehouse_id=?, provider_id=?, status=?, item_name=?, unit_price=?, description=? WHERE id= ?",
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
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: itemUpdate})
	}
}
