package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 中间件：为每个请求生成唯一 TraceID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}

// StructuredLogger 结构化日志中间件
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		traceID, _ := c.Get("trace_id")

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
			slog.Duration("latency", latency),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.Any("trace_id", traceID),
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			slog.LogAttrs(nil, slog.LevelError, "HTTP Request", attrs...)
		case status >= 400:
			slog.LogAttrs(nil, slog.LevelWarn, "HTTP Request", attrs...)
		default:
			slog.LogAttrs(nil, slog.LevelInfo, "HTTP Request", attrs...)
		}
	}
}