package user

import (
	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/user/api"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", api.NewHealthHandler())
	return r
}
