package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupMiddlewareRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCORS(t *testing.T) {
	r := setupMiddlewareRouter()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header not set correctly")
	}
}

func TestRequestID(t *testing.T) {
	r := setupMiddlewareRouter()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		traceID, exists := c.Get("trace_id")
		if !exists {
			t.Error("trace_id not set")
		}
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header missing")
	}
}

func TestMetricsRecord(t *testing.T) {
	AppMetrics = &Metrics{
		RequestsByPath:   make(map[string]int64),
		RequestsByStatus: make(map[int]int64),
		StartTime:        time.Now(),
	}

	AppMetrics.RecordRequest("/api/test", 200, 50)
	AppMetrics.RecordRequest("/api/test", 200, 100)
	AppMetrics.RecordRequest("/api/test", 500, 200)

	if AppMetrics.TotalRequests != 3 {
		t.Errorf("expected 3 requests, got %d", AppMetrics.TotalRequests)
	}
}

func TestRecovery(t *testing.T) {
	r := setupMiddlewareRouter()
	r.Use(Recovery())
	r.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", w.Code, http.StatusInternalServerError)
	}
}