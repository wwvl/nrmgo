package latency_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"nrmgo/internal/latency"
)

func Example() {
	// 创建测试目标
	targets := []latency.Target{
		latency.NewTarget("Google", "https://www.google.com").
			WithTimeout(3 * time.Second),
		latency.NewTarget("GitHub", "https://github.com").
			WithHeaders(http.Header{
				"Accept": []string{"application/json"},
			}),
	}

	// 创建测试器
	tester := latency.NewTester(latency.DefaultOptions().
		WithConcurrency(2).
		WithMaxLatency(5 * time.Second))

	// 执行测试
	ctx := context.Background()
	results := tester.Test(ctx, targets)

	// 处理结果
	for _, result := range results {
		if result.IsOnline {
			fmt.Printf("%s is online, latency: %v\n", result.Name, result.Latency)
		} else {
			fmt.Printf("%s is offline: %s\n", result.Name, result.Error)
		}
	}
}

func ExampleTarget() {
	// 创建带自定义验证的目标
	target := latency.NewTarget("Custom", "https://api.example.com").
		WithTimeout(2 * time.Second).
		WithHeaders(http.Header{
			"Authorization": []string{"Bearer token"},
		}).
		WithValidation(func(r *latency.Result) error {
			if r.Latency > time.Second {
				return fmt.Errorf("latency %v exceeds 1s", r.Latency)
			}
			return nil
		})

	// 创建测试器
	tester := latency.NewTester(nil) // 使用默认选项

	// 执行单个目标测试
	ctx := context.Background()
	result := tester.TestOne(ctx, target)

	// 处理结果
	if result.IsOnline {
		fmt.Printf("Target is online with latency: %v\n", result.Latency)
	} else {
		fmt.Printf("Target is offline: %s\n", result.Error)
	}
}
