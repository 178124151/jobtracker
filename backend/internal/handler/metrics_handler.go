package handler

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics 采集器（轻量级，不依赖 prometheus client）
type Metrics struct {
	mu              sync.RWMutex
	TotalRequests   int64            `json:"total_requests"`
	RequestsByPath  map[string]int64 `json:"requests_by_path"`
	RequestsByStatus map[int]int64   `json:"requests_by_status"`
	AvgResponseMs   float64          `json:"avg_response_ms"`
	StartTime       time.Time        `json:"start_time"`
	LastRequestTime time.Time        `json:"last_request_time"`
}

var AppMetrics = &Metrics{
	RequestsByPath:   make(map[string]int64),
	RequestsByStatus: make(map[int]int64),
	StartTime:        time.Now(),
}

// RecordRequest 记录一次请求
func (m *Metrics) RecordRequest(path string, status int, latencyMs int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++
	m.RequestsByPath[path]++
	m.RequestsByStatus[status]++
	m.LastRequestTime = time.Now()

	// 简单移动平均
	m.AvgResponseMs = (m.AvgResponseMs*float64(m.TotalRequests-1) + float64(latencyMs)) / float64(m.TotalRequests)
}

// MetricsMiddleware 采集请求指标的中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		AppMetrics.RecordRequest(c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

// GetMetrics 返回当前指标
func GetMetrics(c *gin.Context) {
	AppMetrics.mu.RLock()
	defer AppMetrics.mu.RUnlock()

	uptime := time.Since(AppMetrics.StartTime)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"total_requests":   AppMetrics.TotalRequests,
			"requests_by_path": AppMetrics.RequestsByPath,
			"requests_by_status": AppMetrics.RequestsByStatus,
			"avg_response_ms":  AppMetrics.AvgResponseMs,
			"uptime_seconds":   int(uptime.Seconds()),
			"uptime_human":     uptime.String(),
			"start_time":       AppMetrics.StartTime,
			"last_request":     AppMetrics.LastRequestTime,
		},
		"message": "ok",
	})
}