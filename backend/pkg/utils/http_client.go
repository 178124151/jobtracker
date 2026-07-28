package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient 封装的 HTTP 客户端（带超时和重试）
type HTTPClient struct {
	client  *http.Client
	retries int
}

// NewHTTPClient 创建带超时的客户端
func NewHTTPClient(timeout time.Duration, retries int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		retries: retries,
	}
}

// Get 带重试的 GET 请求
func (h *HTTPClient) Get(url string) (int, []byte, error) {
	var lastErr error

	for i := 0; i <= h.retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), h.client.Timeout)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			cancel()
			return 0, nil, fmt.Errorf("create request failed: %w", err)
		}

		resp, err := h.client.Do(req)
		cancel()

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		return resp.StatusCode, body, nil
	}

	return 0, nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}