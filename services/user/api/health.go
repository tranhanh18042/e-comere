package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status int `json:"status"`
}

func NewHealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, HealthResponse{Status: http.StatusOK})
	}
}
