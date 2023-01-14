package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	Status      int    `json:"Status"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
	Email       string `json:"Email"`
}

type CustomerDetail struct {
	Id          string `json:"Id"`
	Status      string `json:"Status"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
	Email       string `json:"Email"`
}

func CreateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Customer Customer
		if err := ctx.ShouldBindJSON(&Customer); err == nil {
			_, err := customerDB.Exec("insert into customer(status, username, password, first_name, last_name, address, phone_number, email) values(?,?,?,?,?,?,?,?)", Customer.Status, Customer.Username, Customer.Password, Customer.FirstName, Customer.LastName, Customer.Address, Customer.PhoneNumber, Customer.Email)
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
		row := customerDB.QueryRow("select * from customer where id = " + ctx.Param("id"))
		if err := row.Scan(&CustomerDetail, &CustomerDetail.Status, &CustomerDetail.Username, &CustomerDetail.Password, &CustomerDetail.FirstName, &CustomerDetail.LastName, &CustomerDetail.Address, &CustomerDetail.PhoneNumber, &CustomerDetail.Email); err == nil {
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
		rows, err := customerDB.Query("select * from customer")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var singleCustomerDetail CustomerDetail
			if err := rows.Scan(&singleCustomerDetail.Id, &singleCustomerDetail.Status, &singleCustomerDetail.Username, &singleCustomerDetail.Password, &singleCustomerDetail.FirstName, &singleCustomerDetail.LastName, &singleCustomerDetail.Address, &singleCustomerDetail.PhoneNumber, &singleCustomerDetail.Email); err == nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message": err.Error(),
				})
				return
			}
			customerDetail = append(customerDetail, singleCustomerDetail)
		}
	}
}
func UpdateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var CustomerDetail CustomerDetail
		if err := ctx.ShouldBindJSON(&CustomerDetail); err == nil {
			update, err := customerDB.Prepare("update customer set status=?, username=?,password=? ,first_name=?, last_name=?, address=?, phone_number=?, email=? where id=" + ctx.Param("id"))
			if err != nil {
				panic(err)
			}
			update.Exec(CustomerDetail.Status, CustomerDetail.Username, CustomerDetail.Password, CustomerDetail.FirstName, CustomerDetail.LastName, CustomerDetail.Address, CustomerDetail.PhoneNumber, CustomerDetail.Email)
			ctx.JSON(200, CustomerDetail)
		} else {
			ctx.JSON(500, gin.H{
				"message": "error",
			})
		}
	}
}
