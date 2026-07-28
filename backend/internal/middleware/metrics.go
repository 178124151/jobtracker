package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	Mu              sync.RWMutex
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

func (m *Metrics) RecordRequest(path string, status int, latencyMs int64) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.TotalRequests++
	m.RequestsByPath[path]++
	m.RequestsByStatus[status]++
	m.LastRequestTime = time.Now()
	m.AvgResponseMs = (m.AvgResponseMs*float64(m.TotalRequests-1) + float64(latencyMs)) / float64(m.TotalRequests)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		AppMetrics.RecordRequest(c.Request.URL.Path, c.Writer.Status(), latency)
	}
}