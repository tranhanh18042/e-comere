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

type CustomerRequest struct {
	Status      int    `json:"status" db:"Status"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Address     string `json:"address" db:"address"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
}

func CreateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customerReq CustomerRequest

		if err := ctx.ShouldBindJSON(&customerReq); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "customer",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()
			logger.Debug(ctx, "invalid params", err)
			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}

		labels := helper.DBMetricsLabels{
			DBName: customerDB.Name,
			Target: "create-customer",
		}
		result, err := helper.DBExecWithMetrics(labels, customerDB,
			"INSERT INTO customer(status, username, password, first_name, last_name, address, phone_number, email) VALUES(?,?,?,?,?,?,?,?)",
			customerReq.Status,
			customerReq.Username,
			customerReq.Password,
			customerReq.FirstName,
			customerReq.LastName,
			customerReq.Address,
			customerReq.PhoneNumber,
			customerReq.Email)

		if err != nil {
			logger.Debug(ctx, "internal error", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		customerID, err := result.LastInsertId()
		if err != nil {
			logger.Debug(ctx, "internal error", err)
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}

		createdCustomer := model.Customer{
			Id: int(customerID),
			Username: customerReq.Username,
			FirstName: customerReq.FirstName,
			LastName: customerReq.LastName,
			Address: customerReq.Address,
			PhoneNumber: customerReq.PhoneNumber,
			Email: customerReq.Email,
		}
		logger.Debug(ctx, "created customer", createdCustomer)
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: createdCustomer})
	}
}

func GetCustomerByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customer model.Customer
		labels := helper.DBMetricsLabels{
			DBName: customerDB.Name,
			Target: "get-customer-by-id",
		}
		err := helper.DBGetWithMetrics(labels, customerDB, &customer, "select * from customer where id = ?", ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse{
			Payload: customer,
		})
	}
}

func GetListCustomers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customers []model.Customer
		labels := helper.DBMetricsLabels{
			DBName: customerDB.Name,
			Target: "get-list-customers",
		}
		err := helper.DBSelectWithMetrics(labels, customerDB, &customers, "select * from customer")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Error(ctx, "item svc get-list-customers internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}

		ctx.JSON(200, helper.SuccessResponse{
			Payload: customers,
		})
	}
}

func UpdateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customerUpdate model.Customer
		if err := ctx.ShouldBindJSON(&customerUpdate); err != nil {
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "customer",
				"path": ctx.FullPath(),
				"type": helper.MetricInvalidParams,
				"env":  "local",
			}).Inc()

			ctx.JSON(http.StatusBadRequest, helper.BadRequestResponse)
			return
		}

		labels := helper.DBMetricsLabels{
			DBName: customerDB.Name,
			Target: "update-item",
		}

		_, err := helper.DBExecWithMetrics(labels, customerDB,
			"UPDATE customer SET status=?, username=?,password=? ,first_name=?, last_name=?, address=?, phone_number=?, email=? WHERE id=?",
			customerUpdate.Status,
			customerUpdate.Username,
			customerUpdate.Password,
			customerUpdate.FirstName,
			customerUpdate.LastName,
			customerUpdate.Address,
			customerUpdate.PhoneNumber,
			customerUpdate.Email,
			ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(200, customerUpdate)
	}
}
