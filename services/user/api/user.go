package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type user struct {
	Id           int    `json:"id"`
	Role_id      int    `json:"role_id`
	Status       int    `json:"status"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Phone_number string `json:"phone_number"`
}
type role struct {
	Id        int    `json:"id"`
	Role_name string `json:"role_name"`
}
type userDTO struct {
	Role_id      role
	Status       int    `json:"status"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Phone_number string `json:"phone_number"`
}

func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var user userDTO
		if err := ctx.ShouldBindJSON(&user); err == nil {
			_, err = db.Exec("insert into user(role_id,status,username,password,name,address,phone_number) values("+ctx.Param("role_id")+",?,?,?,?,?,?)", user.Status, user.Username, user.Password, user.Name, user.Address, user.Phone_number)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err,
				})
				return
			}
			ctx.JSON(200, user)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		}
	}
}
