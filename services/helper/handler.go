package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ToGinHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
