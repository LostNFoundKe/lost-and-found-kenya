package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// RequestLogger middleware for logging HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// Access the status we are sending
		status := c.Writer.Status()

		// Log request details
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("[%s] %s %s %s %s %s %s status: %d\n",
			end.Format("2006-01-02 15:04:05"),
			method,
			path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Errors.String(),
			latency.String(),
			status)))
	}
}
