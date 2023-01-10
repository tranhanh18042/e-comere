package user

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/middlewares"
	"github.com/tranhanh18042/e-comere/services/user/api"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.NewMetricsMiddleware(helper.MetricSvcNameUser))
	r.GET("/api/health", api.NewHealthHandler())
	r.GET("/api/user/info", api.NewUserInfoHandler())
	r.GET("/metrics", helper.ToGinHandler(promhttp.Handler()))
	return r
}
