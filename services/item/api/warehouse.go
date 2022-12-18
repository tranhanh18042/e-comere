package api

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

type Warehouse struct {
	Name_warehouse string `json:"name_warehouse"`
	Address        string `json:"address"`
	Phone_Number   string `json:"phone_number"`
}
type WarehouseId struct {
	Id             int    `json:"id"`
	Name_warehouse string `json:"name_warehouse"`
	Address        string `json:"address"`
	Phone_Number   string `json:"phone_number"`
}

var dbItem *sql.DB

func CreatWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbItem, errItem := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if errItem != nil {
			panic(errItem)
		}
		var warehouse Warehouse

		if err := ctx.ShouldBindJSON(&warehouse); err == nil {
			_, err := dbItem.Exec("INSERT INTO warehouse(name_warehouse, address, phone_number) VALUES(?,?,?)", warehouse.Name_warehouse, warehouse.Address, warehouse.Phone_Number)
			if err != nil {
				ctx.JSON(500, gin.H{
					"messages": err,
				})
			}

			ctx.JSON(200, warehouse)

		} else {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

func GetWarehouseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbItem, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		row := dbItem.QueryRow("SELECT id, name_warehouse, address,phone_number FROM warehouse WHERE id= " + ctx.Param("id"))
		var warehouseId WarehouseId
		if err := row.Scan(&warehouseId.Id, &warehouseId.Name_warehouse, &warehouseId.Address, &warehouseId.Phone_Number); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"messages": "error ",
			})
		}
		ctx.JSON(200, warehouseId)

	}
}

func GetWarehouseAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbItem, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		rows, err := dbItem.Query("SELECT * FROM warehouse")
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "error",
			})
		}
		var warehouseId []WarehouseId

		for rows.Next() {
			var singleWarehouse WarehouseId
			if err := rows.Scan(&singleWarehouse.Id, &singleWarehouse.Name_warehouse, &singleWarehouse.Address, &singleWarehouse.Phone_Number); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"messages": "error 0 ",
				})
			}
			warehouseId = append(warehouseId, singleWarehouse)
		}
		ctx.JSON(200, warehouseId)
	}
}

func UpdateWarehouse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbItem, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var warehouse Warehouse
		if err := ctx.ShouldBindJSON(&warehouse); err == nil {
			update, err := dbItem.Prepare("UPDATE warehouse SET name_warehouse=?, address=?, phone_number=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err.Error())
			}
			update.Exec(warehouse.Name_warehouse, warehouse.Address, warehouse.Phone_Number)

			ctx.JSON(200, warehouse)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}
	}
}
