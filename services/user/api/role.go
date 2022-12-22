package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Role struct {
	Role_name string `json:"role_name"`
}
type Roles struct {
	Id        int    `json:"id"`
	Role_name string `json:"role_name"`
}

var db *sqlx.DB

func CreateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var role Role
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		if err = ctx.ShouldBindJSON(&role); err == nil {
			_, err := db.Exec("INSERT INTO role(role_name) VALUES(?)", role.Role_name)
			if err != nil {
				ctx.JSON(500, gin.H{"message": "error inserting role"})
				return
			}
			ctx.JSON(200, role)
		} else {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}
	}
}
func GetRoleAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		rows, err := db.Query("select * from role")
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err,
			})
		}
		var roles []Roles
		for rows.Next() {
			var singleRoles Roles
			if err := rows.Scan(&singleRoles.Id, &singleRoles.Role_name); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"messages": "error 0 ",
				})
			}
			roles = append(roles, singleRoles)
		}
		ctx.JSON(200, roles)
	}
}

func UpdateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sqlx.Connect("mysql", "root:root@tcp(db_ecom_user:3306)/ecom_user?collation=utf8mb4_unicode_ci&parseTime=true")
		if err != nil {
			panic(err)
		}
		var roles Role
		if err := ctx.ShouldBindJSON(&roles); err == nil {
			update, err := db.Prepare("UPDATE role SET role_name=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err.Error())
			}
			update.Exec(roles.Role_name)

			ctx.JSON(200, roles)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}
	}
}
