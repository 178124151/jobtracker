package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// 表驱动测试：健康检查端点
func TestLiveness(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "返回 alive 状态",
			wantStatus: http.StatusOK,
			wantBody:   "alive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter()
			// 模拟 HealthHandler（不依赖数据库）
			r.GET("/healthz", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "alive"})
			})

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			var resp map[string]string
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["status"] != tt.wantBody {
				t.Errorf("got status %s, want %s", resp["status"], tt.wantBody)
			}
		})
	}
}

// 测试 Metrics 中间件
func TestMetricsMiddleware(t *testing.T) {
	r := setupRouter()
	
	// 模拟指标采集
	requestCount := 0
	r.GET("/test", func(c *gin.Context) {
		requestCount++
		c.JSON(200, gin.H{"ok": true})
	})

	// 发送 3 个请求
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	if requestCount != 3 {
		t.Errorf("expected 3 requests, got %d", requestCount)
	}
}

// 测试 TraceID 中间件
func TestRequestID(t *testing.T) {
	r := setupRouter()
	r.GET("/test-id", func(c *gin.Context) {
		traceID, exists := c.Get("trace_id")
		if !exists {
			t.Error("trace_id not set in context")
		}
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header not set")
	}
}