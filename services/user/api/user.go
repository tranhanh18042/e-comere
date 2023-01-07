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
type info_user struct {
	Id           int    `json:"id"`
	Role_id      int    `json:"role_id`
	Status       int    `json:"status"`
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

var dbUser *sqlx.DB

func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbUser, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var user userDTO
		if err := ctx.ShouldBindJSON(&user); err == nil {
			_, err = dbUser.Exec("insert into user(role_id,status,username,password,name,address,phone_number) values("+ctx.Param("role_id")+",?,?,?,?,?,?)", user.Status, user.Username, user.Password, user.Name, user.Address, user.Phone_number)
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
func GetAllUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbUser, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		rows, err := dbUser.Query("select * from user")
		if err != nil {
			panic(err)
		}
		var users []user
		for rows.Next() {
			var user user
			if err := rows.Scan(&user.Id, &user.Role_id, &user.Status, &user.Name, &user.Address, user.Phone_number); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err,
				})
				return
			}
			users = append(users, user)
		}
		ctx.JSON(200, users)
	}
}
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbUser, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var user user
		row := dbUser.QueryRow("select * from user where id = " + ctx.Param(("id")))
		if err := row.Scan(&user.Id, &user.Role_id, &user.Status, &user.Name, &user.Address, user.Phone_number); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err,
			})
			return
		}
		ctx.JSON(200, user)
	}
}
func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbUser, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var user user
		if err := ctx.ShouldBindJSON(&user); err == nil {
			update, err := dbUser.Preparex("update user set role_id=?,status=?,name=?,addres=?,phone_number=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(&user.Role_id, &user.Status, &user.Name, &user.Address, user.Phone_number)
			ctx.JSON(200, user)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}
	}
}
