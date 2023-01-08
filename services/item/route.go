package item

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tranhanh18042/e-comere/services/item/api"
)

func toGinHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func InitRoute() *gin.Engine {
	r := gin.Default()
	// r.GET("/api/health", api.NewHealthHandler())
	r.POST("/api/warehouse", api.CreatWarehouse())
	r.GET("api/warehouse/:id", api.GetWarehouseById())
	r.GET("/api/warehouse", api.GetWarehouseAll())
	r.PUT("api/warehouse/:id", api.UpdateWarehouse())

	r.POST("/api/provider", api.CreatProvider())
	r.GET("/api/provider/:id", api.GetProviderId())
	r.GET("api/provider", api.GetProviderAll())
	r.PUT("/api/provider/:id", api.UpdateProvider())

	r.POST("/api/item/:warehouse_id/:provider_id", api.AddItem())
	r.GET("api/item", api.GetAllItem())
	r.GET("api/item/:id", api.GetItem())
	r.PUT("api/item/:id", api.UpdateItem())
	r.GET("/metrics", toGinHandler(promhttp.Handler()))
	return r
}
