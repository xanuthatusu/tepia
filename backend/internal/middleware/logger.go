package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		gin.DefaultWriter.Write([]byte(
			time.Now().Format("2006-01-02 15:04:05") +
				" | " + method + " " + path +
				" | " + latency.String() +
				" | Status: " + httpStatusText(status) + "\n",
		))
	}
}

func httpStatusText(code int) string {
	return http.StatusText(code)
}
