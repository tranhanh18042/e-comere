package user

import (
	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/user/api"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", api.NewHealthHandler())
	r.POST("/api/role", api.CreateRole())
	r.GET("/api/role", api.GetRoleAll())
	r.PUT("/api/role/:id", api.UpdateRole())
	r.POST("api/user/:role_id", api.CreateUser())
	return r
}
