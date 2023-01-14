package api

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type Provider struct {
	Name_provider string `json:"name_provider"`
	Phone_number  string `json:"phone_number"`
	Address       string `json:"address"`
}

type ProviderID struct {
	Id            int    `json:"id"`
	Name_provider string `json:"name_provider"`
	Phone_number  string `json:"phone_number"`
	Address       string `json:"address"`
}

func CreatProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var provider Provider
		if err := ctx.ShouldBindJSON(&provider); err == nil {
			_, err := itemDB.Exec("INSERT INTO provider(name_provider,phone_number,address) VALUES(?,?,?)", provider.Name_provider, provider.Phone_number, provider.Address)
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
			ctx.JSON(200, provider)
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
func GetProviderId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerId ProviderID
		row := itemDB.QueryRow("select * from provider where id = " + ctx.Param(("id")))
		if err := row.Scan(&providerId.Id, &providerId.Name_provider, &providerId.Phone_number, providerId.Address); err == nil {
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
		ctx.JSON(200, providerId)
	}
}
func GetProviderAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := itemDB.Query("select * from provider")
		if err != nil {
			panic(err)
		}
		var providerId []ProviderID
		for rows.Next() {
			var singleProviderId ProviderID
			if err := rows.Scan(&singleProviderId.Id, &singleProviderId.Name_provider, &singleProviderId.Phone_number, singleProviderId.Address); err == nil {
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
			providerId = append(providerId, singleProviderId)
		}
		ctx.JSON(200, providerId)
	}
}
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var provider Provider
		if err := ctx.ShouldBindJSON(&provider); err == nil {
			update, err := itemDB.Preparex("update provider set name_provider=?,phone_number=?,address=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(provider.Name_provider, provider.Phone_number, provider.Address)
			ctx.JSON(200, provider)
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
