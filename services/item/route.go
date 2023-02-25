package item

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/item/api"
	"github.com/tranhanh18042/e-comere/services/middlewares"
)

func InitRoute(db *helper.SvcDB) *gin.Engine {
	api.UseDB(db)

	r := gin.Default()
	r.Use(middlewares.NewMetricsMiddleware(helper.MetricSvcNameItem))
	r.GET("/api/health", api.NewHealthHandler())

	r.POST("/api/warehouse", api.CreatWarehouse())
	r.GET("/api/warehouse/:id", api.GetWarehouseById())
	r.GET("/api/warehouse", api.GetListWarehouse())
	r.PUT("/api/warehouse/:id", api.UpdateWarehouse())

	r.POST("/api/provider", api.CreatProvider())
	r.GET("/api/provider/:id", api.GetProviderById())
	r.GET("/api/provider", api.GetListProviders())
	r.PUT("/api/provider/:id", api.UpdateProvider())

	r.POST("/api/item", api.CreateItem())
	r.GET("/api/item", api.GetListItems())
	r.GET("/api/item/:id", api.GetItemByID())
	r.PUT("/api/item/:id", api.UpdateItem())

	r.GET("/metrics", helper.ToGinHandler(promhttp.Handler()))
	return r
}
