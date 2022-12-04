package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranhanh18042/e-comere/services/database"
)

type role struct {
	id        int    `json:"id"`
	role_name string `json:"role_name"`
}
type addRole struct {
	role_name string `json:"role_name"`
}

func CreateRole(c *gin.Context) {

	db := database.DBUserConn()
	var requestBody addRole

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error creating role",
		})
		return
	}
	res, err := db.Exec("insert into role (role_name) values(?);",
		requestBody.role_name,
	)
	id, _ := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"id":        id,
			"role_name": requestBody.role_name,
		})
	}
	defer db.Close()
}

func GetListRoles(g *gin.Context) {

	var roles []role
	db := database.DBUserConn()
	err := db.Select(roles, "select * from role")
	if err != nil {
		g.JSON(500, gin.H{
			"message": "error querying role",
		})
	}

}
