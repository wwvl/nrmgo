package latency

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"nrmgo/internal/config"
)

// Tester 定义延迟测试器接口
type Tester interface {
	// Test 测试多个目标的延迟
	Test(ctx context.Context, targets []Target) []*Result
	// TestOne 测试单个目标的延迟
	TestOne(ctx context.Context, target Target) *Result
}

// DefaultTester 默认的延迟测试器实现
type DefaultTester struct {
	opts   *Options
	client *http.Client
}

// NewTester 创建新的延迟测试器
func NewTester(opts *Options) Tester {
	if opts == nil {
		opts = DefaultOptions()
	}

	client := &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 10,
		},
	}

	return &DefaultTester{
		opts:   opts,
		client: client,
	}
}

// NewTesterFromConfig 从配置创建新的延迟测试器
func NewTesterFromConfig(cfg *config.Config) Tester {
	opts := DefaultOptions()
	if cfg != nil && cfg.MaxConcurrentRequests > 0 {
		opts.Concurrency = cfg.MaxConcurrentRequests
	}
	return NewTester(opts)
}

// Test 并发测试多个目标的延迟
func (t *DefaultTester) Test(ctx context.Context, targets []Target) []*Result {
	if len(targets) == 0 {
		return nil
	}

	results := make([]*Result, len(targets))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, t.opts.Concurrency)

	for i, target := range targets {
		wg.Add(1)
		go func(index int, tgt Target) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			results[index] = t.TestOne(ctx, tgt)
		}(i, target)
	}

	wg.Wait()
	return results
}

// TestOne 测试单个目标的延迟
func (t *DefaultTester) TestOne(ctx context.Context, target Target) *Result {
	result := &Result{
		Name:     target.Name,
		URL:      target.URL,
		TestTime: time.Now(),
	}

	// 使用目标特定的超时或默认超时
	timeout := t.opts.Timeout
	if target.Timeout > 0 {
		timeout = target.Timeout
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 获取完整的测试 URL
	testURL := target.GetTestURL()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result
	}

	// 设置请求头
	req.Header.Set("User-Agent", t.opts.UserAgent)
	for k, values := range target.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	// 执行请求并测量延迟
	start := time.Now()
	resp, err := t.client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Error = fmt.Sprintf("request failed: %v", err)
		result.IsOnline = false
		return result
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		result.IsOnline = false
		return result
	}

	// 检查延迟是否超过最大限制
	if t.opts.MaxLatency > 0 && result.Latency > t.opts.MaxLatency {
		result.Error = fmt.Sprintf("latency too high: %v > %v", result.Latency, t.opts.MaxLatency)
		result.IsOnline = false
		return result
	}

	// 执行自定义验证
	if target.Validate != nil {
		if err := target.Validate(result); err != nil {
			result.Error = fmt.Sprintf("validation failed: %v", err)
			result.IsOnline = false
			return result
		}
	}

	result.IsOnline = true
	return result
}
