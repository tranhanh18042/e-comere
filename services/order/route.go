package order

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/middlewares"
	"github.com/tranhanh18042/e-comere/services/order/api"
)

func InitRoute(db *helper.SvcDB) *gin.Engine {
	api.UseDB(db)

	r := gin.Default()
	r.Use(middlewares.NewMetricsMiddleware(helper.MetricSvcNameOrder))
	r.GET("/api/health", api.NewHealthHandler())
	r.GET("/metrics", helper.ToGinHandler(promhttp.Handler()))

	r.POST("/api/order", api.CreateOrder())
	r.PUT("/api/order/:id", api.UpdateOrder())
	r.GET("/api/order", api.GetListOrders())
	r.GET("/api/order/:id", api.GetOrderByID())
	return r
}
