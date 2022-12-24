package api

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	Name_provider string `json:"name_provider"`
	Phone_number  string `json:"phone_number"`
	Address       string `json:"address"`
}

type ProviderID struct {
	Id            int    `json:"id"`
	Name_provider string `json:"name_provider"`
	Phone_number  string `json:"phone_number"`
	Address       string `json:"address"`
}

var db *sqlx.DB

func CreatProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var provider Provider
		if err := ctx.ShouldBindJSON(&provider); err == nil {
			_, err := db.Exec("INSERT INTO provider(name_provider,phone_number,address) VALUES(?,?,?)", provider.Name_provider, provider.Phone_number, provider.Address)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err,
				})
				return
			}
			ctx.JSON(200, provider)
		} else {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}
	}
}
func GetProviderId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var providerId ProviderID
		row := db.QueryRow("select * from provider where id = " + ctx.Param(("id")))
		if err := row.Scan(&providerId.Id, &providerId.Name_provider, providerId.Phone_number, providerId.Address); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err,
			})
			return
		}
		ctx.JSON(200, providerId)
	}
}
func GetProviderAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		rows, err := db.Query("select * from provider")
		if err != nil {
			panic(err)
		}
		var providerId []ProviderID
		for rows.Next() {
			var singleProviderId ProviderID
			if err := rows.Scan(&singleProviderId.Id, &singleProviderId.Name_provider, &singleProviderId.Phone_number, singleProviderId.Address); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err,
				})
			}
			providerId = append(providerId, singleProviderId)
		}
		ctx.JSON(200, providerId)
	}
}
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_item:3306)/ecom_item?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var provider Provider
		if err := ctx.ShouldBindJSON(&provider); err == nil {
			update, err := db.Preparex("update provider set name_provider=?,phone_number=?,address=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(provider.Name_provider, provider.Phone_number, provider.Address)
			ctx.JSON(200, provider)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}

	}
}
