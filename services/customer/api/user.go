package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerInfoResponse struct {
	CustomerName string `json:"customer_name"`
}

func NewCustomerInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, CustomerInfoResponse{CustomerName: "Hanhtran"})
	}
}
