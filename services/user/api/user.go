package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfoResponse struct {
	Username string `json:"username"`
}

func NewUserInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, UserInfoResponse{Username: "Hanhtran"})
	}
}
