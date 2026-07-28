package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jobtracker/backend/internal/middleware"
)

func GetMetrics(c *gin.Context) {
	m := middleware.AppMetrics
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	uptime := time.Since(m.StartTime)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"total_requests":     m.TotalRequests,
			"requests_by_path":   m.RequestsByPath,
			"requests_by_status": m.RequestsByStatus,
			"avg_response_ms":    m.AvgResponseMs,
			"uptime_seconds":     int(uptime.Seconds()),
			"uptime_human":       uptime.String(),
			"start_time":         m.StartTime,
			"last_request":       m.LastRequestTime,
		},
		"message": "ok",
	})
}