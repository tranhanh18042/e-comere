package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	Status      int    `json:"status" db:"Status"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Address     string `json:"address" db:"address"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
}

type CustomerDetail struct {
	Id          int    `json:"id" db:"id"`
	Status      int    `json:"status" db:"status"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Address     string `json:"address" db:"address"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
}

func CreateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Customer Customer
		if err := ctx.ShouldBindJSON(&Customer); err == nil {
			_, err := customerDB.Exec("INSERT INTO customer(status, username, password, first_name, last_name, address, phone_number, email) VALUES(?,?,?,?,?,?,?,?)", Customer.Status, Customer.Username, Customer.Password, Customer.FirstName, Customer.LastName, Customer.Address, Customer.PhoneNumber, Customer.Email)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err.Error(),
				})
			}
			ctx.JSON(200, Customer)
		} else {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
		}
	}
}
func GetCustomerByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var CustomerDetail CustomerDetail
		row := customerDB.QueryRow("SELECT id, status, username, password, first_name, last_name, address, phone_number, email from customer WHERE id = " + ctx.Param("id"))
		if err := row.Scan(&CustomerDetail.Id, &CustomerDetail.Status, &CustomerDetail.Username, &CustomerDetail.Password, &CustomerDetail.FirstName, &CustomerDetail.LastName, &CustomerDetail.Address, &CustomerDetail.PhoneNumber, CustomerDetail.Email); err == nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, CustomerDetail)
	}
}
func GetListCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customerDetail []CustomerDetail
		rows, err := customerDB.Query("SELECT * FROM customer")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var singleCustomerDetail CustomerDetail
			if err := rows.Scan(&singleCustomerDetail.Id, &singleCustomerDetail.Status, &singleCustomerDetail.Username, &singleCustomerDetail.Password, &singleCustomerDetail.FirstName, &singleCustomerDetail.LastName, &singleCustomerDetail.Address, &singleCustomerDetail.PhoneNumber, singleCustomerDetail.Email); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err,
				})
				return
			}
			customerDetail = append(customerDetail, singleCustomerDetail)
		}
		ctx.JSON(http.StatusOK, customerDetail)
	}
}
func UpdateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var CustomerUpdate Customer
		if err := ctx.ShouldBindJSON(&CustomerUpdate); err == nil {
			update, err := customerDB.Prepare("UPDATE customer SET status=?, username=?,password=? ,first_name=?, last_name=?, address=?, phone_number=?, email=? WHERE id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(CustomerUpdate.Status, CustomerUpdate.Username, CustomerUpdate.Password, CustomerUpdate.FirstName, CustomerUpdate.LastName, CustomerUpdate.Address, CustomerUpdate.PhoneNumber, CustomerUpdate.Email)
			ctx.JSON(200, CustomerUpdate)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}
	}
}
