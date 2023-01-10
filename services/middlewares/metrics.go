package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

func NewMetricsMiddleware(svc string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().UnixMilli()

		ctx.Next()

		end := time.Now().UnixMilli()
		labelsCnt := prometheus.Labels{
			"svc": svc,
			"method": ctx.Request.Method,
			"path":   ctx.FullPath(),
			"env":    "local",
			"status": strconv.Itoa(ctx.Writer.Status()),
		}
		labelsDur := prometheus.Labels{
			"svc": svc,
			"method": ctx.Request.Method,
			"path":   ctx.FullPath(),
			"env":    "local",
		}
		metrics.API.ReqCnt.With(labelsCnt).Inc()
		metrics.API.ReqDur.With(labelsDur).Observe(float64(end - start))
	}
}
