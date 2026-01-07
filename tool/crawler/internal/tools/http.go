package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPConfig HTTP客户端配置
type HTTPConfig struct {
	Timeout       time.Duration
	ProxyURL      string
	UserAgent     string
	Headers       map[string]string
	MaxRetryCount int
}

// HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
	config HTTPConfig
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(config HTTPConfig) *HTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.UserAgent == "" {
		config.UserAgent = GetRandomUserAgent()
	}

	if config.MaxRetryCount == 0 {
		config.MaxRetryCount = 3
	}

	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &HTTPClient{
		client: client,
		config: config,
	}
}

// Get 发送GET请求
func (c *HTTPClient) Get(url string) ([]byte, error) {
	return c.Request(http.MethodGet, url, nil, nil)
}

// Post 发送POST请求
func (c *HTTPClient) Post(url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	return c.Request(http.MethodPost, url, bodyReader, headers)
}

// Request 发送HTTP请求
func (c *HTTPClient) Request(method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	var respBody []byte
	var err error

	for i := 0; i <= c.config.MaxRetryCount; i++ {
		respBody, err = c.doRequest(method, url, body, headers)
		if err == nil {
			return respBody, nil
		}

		// 重试等待
		if i < c.config.MaxRetryCount {
			SleepRandom(1, 3)
		}
	}

	return nil, err
}

// doRequest 执行单次HTTP请求
func (c *HTTPClient) doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置配置的headers
	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}

	// 设置请求的headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
			URL:        url,
		}
	}

	return respBody, nil
}

// Config 获取HTTP客户端配置
func (c *HTTPClient) Config() HTTPConfig {
	return c.config
}

// HTTPError HTTP错误
type HTTPError struct {
	StatusCode int
	Body       string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(`http request failed: url="%s", status_code=%d, body="%s"`, e.URL, e.StatusCode, e.Body)
}
