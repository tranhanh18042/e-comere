package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
)

var customerClient *http.Client
var customerBaseURL = "http://svc_customer:8080"

func init() {
	customerClient = http.DefaultClient
}

func CreateCustomer(customer *model.Customer) (customerId int, err error) {
	logger.Debug(context.Background(), "customer info", *customer)
	data := url.Values{
        "first_name": {customer.FirstName},
		"last_name": {customer.LastName},
		"address": {customer.Address},
		"phone_number": {customer.PhoneNumber},
		"email": {customer.Email},
    }
    resp, err := http.PostForm(customerBaseURL + "/api/customer", data)
    if err != nil {
        return 0, fmt.Errorf("cannot create customer, err: %v", err)
    }

    var res model.Customer
    json.NewDecoder(resp.Body).Decode(&res)
    logger.Debug(context.Background(), "create customer", res)

	return res.Id, nil
}

func GetCustomerByID(customerID int) (*model.Customer, error) {
	resp, err := customerClient.Get(customerBaseURL + "/api/customer/"+strconv.FormatInt(int64(customerID), 10))
	if err != nil {
		return nil,  fmt.Errorf("cannot get customer with id: %d, err: %v", customerID, err)
	}
	var customer model.Customer
	json.NewDecoder(resp.Body).Decode(&customer)
    logger.Debug(context.Background(), "get customer", customer)

	return &customer, nil
}
