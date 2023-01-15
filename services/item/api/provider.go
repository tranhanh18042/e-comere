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

type ProviderRequest struct {
	ProviderName string `json:"provider_name"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
}

func CreatProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerReq ProviderRequest
		if err := ctx.ShouldBindJSON(&providerReq); err == nil {
			_, err := itemDB.Exec("INSERT INTO provider(provider_name,phone_number,address) VALUES(?,?,?)",
			providerReq.ProviderName,
			providerReq.PhoneNumber,
			providerReq.Address)

			if err != nil {
				metrics.API.ErrCnt.With(prometheus.Labels{
					"svc":  "item",
					"path": ctx.FullPath(),
					"type": helper.MetricInvalidParams,
					"env":  "local",
				}).Inc()
				ctx.JSON(500, gin.H{
					"message": err,
				})
				return
			}
			ctx.JSON(200, providerReq)
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
func GetProviderById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var provider model.Provider
		row := itemDB.QueryRow("select * from provider where id = " + ctx.Param(("id")))
		if err := row.Scan(&provider.Id, &provider.ProviderName, &provider.PhoneNumber, provider.Address); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err,
			})
			metrics.API.ErrCnt.With(prometheus.Labels{
				"svc":  "item",
				"path": ctx.FullPath(),
				"type": helper.MetricQueryError,
				"env":  "local",
			}).Inc()
			return
		}
		ctx.JSON(200, provider)
	}
}
func GetProviderAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := itemDB.Query("select * from provider")
		if err != nil {
			panic(err)
		}
		var providers []model.Provider
		for rows.Next() {
			var singleProviderId model.Provider
			if err := rows.Scan(&singleProviderId.Id, &singleProviderId.ProviderName, &singleProviderId.PhoneNumber, singleProviderId.Address); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err,
				})
				metrics.API.ErrCnt.With(prometheus.Labels{
					"svc":  "item",
					"path": ctx.FullPath(),
					"type": helper.MetricQueryError,
					"env":  "local",
				}).Inc()
			}
			providers = append(providers, singleProviderId)
		}
		ctx.JSON(200, providers)
	}
}
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerReq ProviderRequest
		if err := ctx.ShouldBindJSON(&providerReq); err == nil {
			update, err := itemDB.Preparex("update provider set provider_name=?,phone_number=?,address=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(providerReq.ProviderName, providerReq.PhoneNumber, providerReq.Address)
			ctx.JSON(200, providerReq)
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
