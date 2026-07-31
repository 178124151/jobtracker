package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/middleware"
)

// GetPrometheusMetrics returns Prometheus format metrics
func GetPrometheusMetrics(c *gin.Context) {
	m := middleware.AppMetrics
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	uptime := time.Since(m.StartTime)

	output := ""
	output += "# HELP http_requests_total Total number of HTTP requests\n"
	output += "# TYPE http_requests_total counter\n"
	output += fmt.Sprintf("http_requests_total %d\n", m.TotalRequests)

	output += "# HELP http_request_duration_ms Average HTTP request duration in milliseconds\n"
	output += "# TYPE http_request_duration_ms gauge\n"
	output += fmt.Sprintf("http_request_duration_ms %.2f\n", m.AvgResponseMs)

	output += "# HELP app_uptime_seconds Application uptime in seconds\n"
	output += "# TYPE app_uptime_seconds gauge\n"
	output += fmt.Sprintf("app_uptime_seconds %d\n", int(uptime.Seconds()))

	output += "# HELP http_requests_by_status_total HTTP requests by status code\n"
	output += "# TYPE http_requests_by_status_total counter\n"
	for status, count := range m.RequestsByStatus {
		output += fmt.Sprintf("http_requests_by_status_total{status=\"%d\"} %d\n", status, count)
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, output)
}