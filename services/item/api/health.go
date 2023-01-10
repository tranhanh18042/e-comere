package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type HealthResponse struct {
	Status int `json:"status"`
}

func NewHealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		metrics.API.ErrCnt.With(prometheus.Labels{
			"svc":  "item",
			"path": ctx.FullPath(),
			"type": ctx.Request.Method,
			"env":  "local",
		}).Inc()
		ctx.JSON(http.StatusBadRequest, HealthResponse{Status: http.StatusBadRequest})
	}
}
