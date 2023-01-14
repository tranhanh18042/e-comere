package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/customer"
	"github.com/tranhanh18042/e-comere/services/item"
	"github.com/tranhanh18042/e-comere/services/order"
)

func main() {
	service, _ := os.LookupEnv("ECOM_SERVICE")

	var router *gin.Engine

	switch service {
	case "svc_item":
		router = item.InitRoute()
	case "svc_customer":
		router = customer.InitRoute()
	case "svc_order":
		router = order.InitRoute()
	default:
		fmt.Println("not support service:", service)
		os.Exit(1)
	}

	router.Run()
}
