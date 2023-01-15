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

type ProviderRequest struct {
	ProviderName string `json:"provider_name"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
}

func GetProviderById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var provider model.Provider
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-provider-by-id",
		}
		err := helper.DBGetWithMetrics(labels, itemDB, &provider, "SELECT * FROM provider WHERE id =? ", ctx.Param("id"))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: provider})
	}
}

func GetListProviders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providers []model.Provider
		labels := helper.DBMetricsLabels{
			DBName: itemDB.Name,
			Target: "get-list-providers",
		}
		err := helper.DBSelectWithMetrics(labels, itemDB, &providers, "SELECT * FROM provider")
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, helper.DataNotFoundResponse)
			} else {
				logger.Error(ctx, "item svc get-list-providers internal error", err)
				ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			}
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{Payload: providers})
	}
}

func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerReq ProviderRequest
		if err := ctx.ShouldBindJSON(&providerReq); err != nil {
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
			Target: "update-provider",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB, "UPDATE provider SET provider_name=?,phone_number=?,address=? WHERE id=?", ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, providerReq)
	}
}

func CreatProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerReq ProviderRequest
		if err := ctx.ShouldBindJSON(&providerReq); err != nil {
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
			Target: "create-provider",
		}
		_, err := helper.DBExecWithMetrics(labels, itemDB, "INSERT INTO provider(provider_name,phone_number,address) VALUES(?,?,?)",
			providerReq.ProviderName,
			providerReq.PhoneNumber,
			providerReq.Address)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.InternalErrorResponse)
			return
		}
		ctx.JSON(http.StatusOK, helper.SuccessResponse{
			Payload: providerReq,
		})
	}
}
