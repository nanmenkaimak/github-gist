package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/gist/metrics"
	"strconv"
	"time"
)

func (h *EndpointHandler) metricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		path := c.Request.URL.Path

		statusString := strconv.Itoa(c.Writer.Status())

		metrics.HttpResponseTime.WithLabelValues(path, statusString, c.Request.Method).Observe(time.Since(start).Seconds())
		metrics.HttpRequestsTotalCollector.WithLabelValues(path, statusString, c.Request.Method).Inc()
	}
}
