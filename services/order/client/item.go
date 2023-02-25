package client

import (
	"context"
	"net/http"
	"strconv"

	"github.com/tranhanh18042/e-comere/services/model"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
)

var itemClient *http.Client

const (
	itemBaseURL = "http://svc_item:8080"
	itemGetPath = "/api/item/:id"
	itemSvc     = "item"
)

func init() {
	itemClient = httpClient
}
func GetItemByID(ItemId int) (*model.Item, error) {
	var item model.Item
	params := &getParams{
		pathParams: map[string]string{"id": strconv.Itoa(ItemId)},
	}
	err := get(itemSvc, itemGetPath, itemClient, itemBaseURL+params.formatURL(itemGetPath), &item)
	if err != nil {
		logger.Error(context.Background(), "get item failed", err)
		return nil, err
	}
	logger.Debug(context.Background(), "get item", item)
	return &item, nil
}
