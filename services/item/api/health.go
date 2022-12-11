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
	return func(ctx *gin.Context)  {
		metrics.API.ReqCnt.With(prometheus.Labels{
				"svc": "item",
				"method": ctx.Request.Method,
				"path":   ctx.FullPath(),
				"env":    "local",
				"status": "200",
			}).Inc()
		ctx.JSON(http.StatusOK, HealthResponse{Status: http.StatusOK})
	}
}
