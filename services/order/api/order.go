package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/order/client"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type OrderRequest struct {
	Status            int    `json:"status"`
	CustomerId        int    `json:"customer_id"`
	CustomerFirstName string `json:"customer_first_name"`
	CustomerLastName  string `json:"customer_last_name"`
	CustomerPhone     string `json:"customer_phone"`
	CustomerEmail     string `json:"customer_email"`
	ItemId            int    `json:"item_id"`
	ItemName          string `json:"item_name"`
	ItemQuantity      int    `json:"item_quantity"`
	ItemAmount        int    `json:"item_amount"`
	Address           string `json:"address"`
	ShipFee           int    `json:"ship_fee"`
	TotalAmount       int    `json:"total_amount"`
	Discount          int    `json:"discount_amount"`
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
			logger.Debug(ctx, "invalid params", err)
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}
		labels := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "create-order",
		}
		if orderReq.CustomerId != 0 {
			customerInfo, err := client.GetCustomerByID(orderReq.CustomerId)
			if err != nil || customerInfo.Id == 0 {
				logger.Debug(ctx, "cannot get customer", err)
				ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
				return
			}
			orderReq.CustomerFirstName = customerInfo.FirstName
			orderReq.CustomerLastName = customerInfo.LastName
			orderReq.CustomerPhone = customerInfo.PhoneNumber
			orderReq.CustomerEmail = customerInfo.Email
		} else {
			customerID, err := client.CreateCustomer(&model.Customer{
				FirstName:   orderReq.CustomerFirstName,
				LastName:    orderReq.CustomerLastName,
				PhoneNumber: orderReq.CustomerPhone,
				Address:     orderReq.Address,
				Email:       orderReq.CustomerEmail,
			})
			if err != nil {
				logger.Debug(ctx, "cannot create customer", err)
				ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
				return
			}
			orderReq.CustomerId = customerID
		}
		if orderReq.ItemId == 0 {
			logger.Debug(ctx, "item does not exist")
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		} else {
			itemInfo, err := client.GetItemByID(orderReq.ItemId)
			if err != nil {
				logger.Debug(ctx, "item does not exist", err)
				ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
				return
			}
			orderReq.Discount = 10
			orderReq.ShipFee = 10
			orderReq.ItemName = itemInfo.ItemName
			orderReq.ItemAmount = itemInfo.UnitPrice
			orderReq.TotalAmount = itemInfo.UnitPrice*orderReq.ItemQuantity + orderReq.ShipFee - orderReq.Discount

		}

		_, err := helper.DBExecWithMetrics(labels, orderDB,
			"INSERT INTO `order`(`status`, customer_id, "+
				"item_id, item_quantity, address, item_amount, "+
				"ship_fee,total_amount, discount_amount) "+
				"VALUES(?,?,?,?,?,?,?,?,?)",
			orderReq.Status,
			orderReq.CustomerId,
			orderReq.ItemId,
			orderReq.ItemQuantity,
			orderReq.Address,
			orderReq.ItemAmount,
			orderReq.ShipFee,
			orderReq.TotalAmount,
			orderReq.Discount)
		if err != nil {
			logger.Debug(ctx, "cannot create order", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: orderReq})
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orderLog model.OrderLog
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
			"UPDATE `order` SET status=?, customer_id=?, item_id=?,"+
				"item_quantity=?, address=?, item_amount=?, ship_fee=?,"+
				"total_amount=?, discount_amount=? WHERE id=?",
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
			logger.Debug(ctx, "cannot update order", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		labelsLog := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "add-event-time",
		}
		orderLog.EventTime = time.Now()
		_, err = helper.DBExecWithMetrics(labelsLog, orderDB, "INSERT INTO `order_log` (order_id,status,event_time) VALUES (?,?,?)",
			ctx.Param("id"),
			orderLog.Status,
			orderLog.EventTime,
		)
		if err != nil {
			logger.Debug(ctx, "cannot insert order_log", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: orderUpdate})
	}
}

func GetListOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Orders []model.Order
		labels := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "get-list-orders",
		}
		err := helper.DBSelectWithMetrics(labels, orderDB, &Orders, "SELECT * FROM `order`")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Debug(ctx, "order svc get-list-order internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: Orders})
	}
}

func GetOrderByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Order model.Order
		labels := helper.DBMetricsLabels{
			DBName: orderDB.Name,
			Target: "get-order-by-id",
		}
		err := helper.DBGetWithMetrics(labels, orderDB, &Order, "SELECT *FROM `order` WHERE id=?", ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Debug(ctx, "order svc get-order-by-id error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: Order})
	}
}
