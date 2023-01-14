package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/customer/api"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/middlewares"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.NewMetricsMiddleware(helper.MetricSvcNameCustomer))
	r.GET("/api/health", api.NewHealthHandler())
	r.GET("/api/customer/info", api.NewCustomerInfoHandler())
	r.GET("/metrics", helper.ToGinHandler(promhttp.Handler()))
	return r
}
