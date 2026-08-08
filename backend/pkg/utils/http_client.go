package utils

import (
	"context"
	"errors"
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

// TimeoutError 表示请求超时或上下文被取消，站点可能仍可访问
type TimeoutError struct {
	Err error
}

func (e *TimeoutError) Error() string { return fmt.Sprintf("timeout: %v", e.Err) }
func (e *TimeoutError) Unwrap() error { return e.Err }

// NewHTTPClient 创建带超时的客户端
func NewHTTPClient(timeout time.Duration, retries int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
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
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := h.client.Do(req)
		cancel()

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		return resp.StatusCode, body, nil
	}

	if lastErr != nil && (errors.Is(lastErr, context.DeadlineExceeded) || errors.Is(lastErr, context.Canceled)) {
		return 0, nil, &TimeoutError{Err: lastErr}
	}
	return 0, nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
