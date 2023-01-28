package client

import (
	"context"
	"net/http"
	"strconv"

	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
)

var customerClient *http.Client
const (
	customerBaseURL = "http://svc_customer:8080"
	customerCreatePath = "/api/customer"
	customerGetPath = "/api/customer/:id"
	customerSvc = "customer"
)

func init() {
	customerClient = httpClient
}

func CreateCustomer(customer *model.Customer) (customerId int, err error) {
	logger.Debug(context.Background(), "customer info", *customer)
    var res model.Customer
    err = post(customerSvc, customerCreatePath, customerClient, customerBaseURL + customerCreatePath, customer, &res)
	if err != nil {
		logger.Error(context.Background(), "create customer err", err)
		return
	}

	logger.Debug(context.Background(), "create customer info", res)
	return res.Id, nil
}

func GetCustomerByID(customerID int) (*model.Customer, error) {
	var customer model.Customer
	params := &getParams{
		pathParams: map[string]string{"id": strconv.Itoa(customerID)},
	}
	err := get(customerSvc, customerGetPath, customerClient, customerBaseURL + params.formatURL(customerGetPath), &customer)
	if err != nil {
		logger.Error(context.Background(), "get customer err", err)
		return nil, err
	}

    logger.Debug(context.Background(), "get customer", customer)
	return &customer, nil
}
