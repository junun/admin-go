package middleware

import (
	"api/models"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tBegin := time.Now()

		// 执行下一步逻辑
		c.Next()

		duration := float64(time.Since(tBegin)) / float64(time.Second)

		// 请求数加1
		models.HTTPReqTotal.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Inc()

		//  记录本次请求处理时间
		models.HTTPReqDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Observe(duration)
	}
}
