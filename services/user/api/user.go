package api

import (
	"github.com/gin-gonic/gin"
)

type user struct {
	Username string `json:"username"`
}

func NewUserInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ctx.JSON(http.StatusOK, UserInfoResponse{Username: "Hanhtran"})
	}
}
