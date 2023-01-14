package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/customer"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/item"
	"github.com/tranhanh18042/e-comere/services/order"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

const (
	svcItem = "svc_item"
	svcOrder = "svc_order"
	svcCustomer = "svc_customer"
)

func main() {
	logger.Init()
	ctx := context.Background()

	logger.Info(ctx, "Init main")

	service, _ := os.LookupEnv("ECOM_SERVICE")

	logger.Info(ctx, "Init DB", service)
	db, err := initDBConn(service)
	if err != nil {
		panic(fmt.Errorf("cannot connect db, err: %v", err))
	}
	defer db.Close()

	logger.Info(ctx, "Init metrics", service)
	go metrics.StatDB("local", "service", db.DB.DB)

	logger.Info(ctx, "Init service", service)
	router, err := initService(service, db)
	if err != nil {
		panic(fmt.Errorf("cannot init service, err: %v", err))
	}

	logger.Info(ctx, "Service start", service)
	router.Run()
}

func initDBConn(service string) (*helper.SvcDB, error) {
	var dbConnStr string
	switch service {
	case svcItem:
		dbConnStr = "root:root@tcp(db_ecom_item:3306)/service_item?collation=utf8mb4_unicode_ci&parseTime=true"
	case svcCustomer:
		dbConnStr = "root:root@tcp(db_ecom_customer:3306)/service_customer?collation=utf8mb4_unicode_ci&parseTime=true"
	case svcOrder:
		dbConnStr = "root:root@tcp(db_ecom_order:3306)/service_order?collation=utf8mb4_unicode_ci&parseTime=true"
	default:
		return nil, fmt.Errorf("invalid service: %s", service)
	}

	return helper.NewDBConn(service, dbConnStr)
}

func initService(service string, db *helper.SvcDB) (*gin.Engine, error) {
	switch service {
	case svcItem:
		return item.InitRoute(db), nil
	case svcCustomer:
		return customer.InitRoute(db), nil
	case svcOrder:
		return order.InitRoute(db), nil
	default:
		return nil, fmt.Errorf("cannot init service for service name: %s", service)
	}
}
