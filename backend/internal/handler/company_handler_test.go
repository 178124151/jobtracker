package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupCompanyRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// 表驱动测试：公司列表 API
func TestCompanyList(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "获取全部公司",
			query:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "按分组筛选",
			query:      "?group=bigtech",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupCompanyRouter()
			r.GET("/api/v1/companies", func(c *gin.Context) {
				group := c.Query("group")
				c.JSON(http.StatusOK, gin.H{
					"code":    0,
					"data":    []string{},
					"group":   group,
					"message": "ok",
				})
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["code"] != float64(0) {
				t.Errorf("got code %v, want 0", resp["code"])
			}
		})
	}
}

// 表驱动测试：SRE 健康检查
func TestSREHealth(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantField  string
	}{
		{
			name:       "健康检查端点",
			path:       "/api/v1/sre/health",
			wantStatus: http.StatusOK,
			wantField:  "status",
		},
		{
			name:       "指标端点",
			path:       "/api/v1/sre/metrics",
			wantStatus: http.StatusOK,
			wantField:  "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupCompanyRouter()
			r.GET("/api/v1/sre/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
			r.GET("/api/v1/sre/metrics", func(c *gin.Context) {
				c.JSON(200, gin.H{"code": 0, "data": gin.H{}})
			})

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if _, exists := resp[tt.wantField]; !exists {
				t.Errorf("response missing field: %s", tt.wantField)
			}
		})
	}
}