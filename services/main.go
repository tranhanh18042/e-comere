package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tranhanh18042/e-comere/services/customer"
	customerDB "github.com/tranhanh18042/e-comere/services/customer/db"
	"github.com/tranhanh18042/e-comere/services/item"
	itemDB "github.com/tranhanh18042/e-comere/services/item/db"
	"github.com/tranhanh18042/e-comere/services/order"
)

func main() {
	service, _ := os.LookupEnv("ECOM_SERVICE")

	var router *gin.Engine
	var db *sqlx.DB
	var errDBConn error

	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	switch service {
	case "svc_item":
		itemDBConn := "root:root@tcp(db_ecom_item:3306)/service_item?collation=utf8mb4_unicode_ci&parseTime=true"
		db, errDBConn = itemDB.NewItemDB(itemDBConn)
		if errDBConn != nil {
			panic(fmt.Errorf("cannot connect to db: %s, err: %v", itemDBConn, errDBConn))
		}
		router = item.InitRoute(db)
	case "svc_customer":
		customerDBConn := "root:root@tcp(db_ecom_customer:3306)/service_customer?collation=utf8mb4_unicode_ci&parseTime=true"
		db, errDBConn = customerDB.NewCustomerDB(customerDBConn)
		if errDBConn != nil {
			panic(fmt.Errorf("cannot connect to db: %s, err: %v", customerDBConn, errDBConn))
		}
		router = customer.InitRoute(db)
	case "svc_order":
		orderDBConn := "root:root@tcp(db_ecom_order:3306)/service_order?collation=utf8mb4_unicode_ci&parseTime=true"
		db, errDBConn = itemDB.NewItemDB(orderDBConn)
		if errDBConn != nil {
			panic(fmt.Errorf("cannot connect to db: %s, err: %v", orderDBConn, errDBConn))
		}
		router = order.InitRoute(db)
	default:
		fmt.Println("not support service:", service)
		os.Exit(1)
	}

	router.Run()
}
