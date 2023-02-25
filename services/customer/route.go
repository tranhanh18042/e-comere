package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/customer/api"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/middlewares"
)

func InitRoute(db *helper.SvcDB) *gin.Engine {
	api.UseDB(db)

	r := gin.Default()
	r.Use(middlewares.NewMetricsMiddleware(helper.MetricSvcNameCustomer))
	r.GET("/api/health", api.NewHealthHandler())

	r.POST("/api/customer", api.CreateCustomer())
	r.PUT("api/customer/:id", api.UpdateCustomer())
	r.GET("/api/customer/:id", api.GetCustomerByID())
	r.GET("/api/customer", api.GetListCustomers())

	r.GET("/metrics", helper.ToGinHandler(promhttp.Handler()))
	return r
}
