package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type OrderRequest struct {
	Status          int    `json:"status"`
	CustomerId      int    `json:"customer_id"`
	ItemId          int    `json:"item_id"`
	ItemQuantity    int    `json:"item_quantity"`
	Address         string `json:"address"`
	ItemAmount      int    `json:"item_account_id"`
	ShipFee         int    `json:"ship_fee"`
	TotalAmount     int    `json:"total_amount"`
	DiscoountAmount int    `json:"discoount_amount"`
}

func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orderReq OrderRequest
		if err := ctx.ShouldBindJSON(&orderReq); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "order",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}
		labels := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "create-order",
		}
		_, err := helper.DBExecWithMetrics(labels, orderDB,
			"INSER INTO order(status, customer_id, item_id, item_quantity, address, item_amount, ship_fee,total_amount, discount_amount) VALUES(?,?,?,?,?,?,?,?,?)",
			orderReq.Status,
			orderReq.CustomerId,
			orderReq.ItemId,
			orderReq.ItemQuantity,
			orderReq.Address,
			orderReq.ItemAmount,
			orderReq.ShipFee,
			orderReq.TotalAmount,
			orderReq.DiscoountAmount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: orderReq})
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orderUpdate model.Order
		if err := ctx.ShouldBindJSON(&orderUpdate); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "order",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()

			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}

		labels := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "update-order",
		}

		_, err := helper.DBExecWithMetrics(labels, orderDB,
			"UPDATE order SET status=?, customer_id=?, item_id=?, item_quantity=?, address=?, item_amount=?, ship_fee=?, totel_amount=?, discount_amount=? WHERE id=?",
			orderUpdate.Status,
			orderUpdate.CustomerId,
			orderUpdate.ItemId,
			orderUpdate.ItemQuantity,
			orderUpdate.Address,
			orderUpdate.ItemAmount,
			orderUpdate.ShipFee,
			orderUpdate.TotalAmount,
			orderUpdate.DiscountAmount,
			ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, orderUpdate)
	}
}
