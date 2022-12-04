package user

import (
	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/user/api"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", api.NewHealthHandler())
	r.GET("/api/user/info", api.NewUserInfoHandler())
	r.POST("api/role/add", api.CreateRole)
	return r
}
