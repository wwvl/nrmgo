package latency

import (
	"net/http"
	"strings"
	"time"
)

// Result 表示延迟测试的结果
type Result struct {
	Name     string        // 目标名称
	URL      string        // 测试的 URL
	IsOnline bool          // 是否在线
	Latency  time.Duration // 延迟时间
	Error    string        // 错误信息
	TestTime time.Time     // 测试时间
}

// Target 表示一个测试目标
type Target struct {
	Name     string         // 目标名称
	URL      string         // 目标 URL
	TestPath string         // 测试路径
	Headers  http.Header    // 请求头
	Timeout  time.Duration  // 超时时间
	Validate ValidationFunc // 自定义验证函数
}

// ValidationFunc 定义自定义验证函数类型
type ValidationFunc func(*Result) error

// Options 表示延迟测试的配置选项
type Options struct {
	Timeout     time.Duration // 请求超时时间
	MaxLatency  time.Duration // 最大可接受延迟
	UserAgent   string        // User-Agent 头
	Concurrency int           // 并发测试数量
}

// DefaultOptions 返回默认的测试选项
func DefaultOptions() *Options {
	return &Options{
		Timeout:     5 * time.Second,
		MaxLatency:  3 * time.Second,
		UserAgent:   "NRMG-Latency-Tester/1.0",
		Concurrency: 5,
	}
}

// NewTarget 创建一个新的测试目标
func NewTarget(name, url string) Target {
	return Target{
		Name:    name,
		URL:     url,
		Headers: make(http.Header),
	}
}

// WithTestPath 设置测试路径
func (t Target) WithTestPath(path string) Target {
	t.TestPath = path
	return t
}

// WithTimeout 设置超时时间
func (t Target) WithTimeout(timeout time.Duration) Target {
	t.Timeout = timeout
	return t
}

// WithHeaders 设置请求头
func (t Target) WithHeaders(headers http.Header) Target {
	t.Headers = headers
	return t
}

// WithValidation 设置自定义验证函数
func (t Target) WithValidation(validate ValidationFunc) Target {
	t.Validate = validate
	return t
}

// GetTestURL 获取完整的测试 URL
func (t Target) GetTestURL() string {
	if t.TestPath == "" {
		return t.URL
	}
	if !strings.HasSuffix(t.URL, "/") {
		return t.URL + "/" + strings.TrimPrefix(t.TestPath, "/")
	}
	return t.URL + strings.TrimPrefix(t.TestPath, "/")
}

// WithTimeout 设置全局超时时间
func (o *Options) WithTimeout(timeout time.Duration) *Options {
	o.Timeout = timeout
	return o
}

// WithMaxLatency 设置最大可接受延迟
func (o *Options) WithMaxLatency(maxLatency time.Duration) *Options {
	o.MaxLatency = maxLatency
	return o
}

// WithUserAgent 设置 User-Agent
func (o *Options) WithUserAgent(userAgent string) *Options {
	o.UserAgent = userAgent
	return o
}

// WithConcurrency 设置并发数量
func (o *Options) WithConcurrency(concurrency int) *Options {
	o.Concurrency = concurrency
	return o
}
