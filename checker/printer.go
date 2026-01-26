package checker

import (
	"fmt"
	"strings"
)

func humanReadableError(err error) string {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "context deadline exceeded"):
		return "请求超时（目标长时间未响应，可能端口被防火墙拦截或服务不可达）"
	case strings.Contains(msg, "no such host"):
		return "DNS 解析失败（域名不存在或无法解析）"
	case strings.Contains(msg, "connection refused"):
		return "连接被拒绝（端口未开放或服务未启动）"
	default:
		return msg
	}
}

// PrintResult 将 ScanResult 转换为“用户可读”的输出文本
func PrintResult(result *ScanResult) (level string, title string, lines []string) {
	title = result.TargetName
	lines = []string{}

	// 请求失败
	if result.Error != nil {
		level = "ERROR"
		lines = append(lines,
			fmt.Sprintf("URL: %s", result.StartURL),
			"检测结果: 请求失败",
			fmt.Sprintf("失败原因: %s", result.Error.Error()),
		)
		if result.StatusCode == 0 {
			lines = append(lines, "HTTP 状态码: 无（未收到响应）")
		} else {
			lines = append(lines, fmt.Sprintf("HTTP 状态码: %d", result.StatusCode))
		}
		return
	}

	// 正常请求完成
	lines = append(lines,
		fmt.Sprintf("起始 URL: %s", result.StartURL),
		fmt.Sprintf("最终 URL: %s", result.FinalURL),
		fmt.Sprintf("页面标题: %s", result.FinalTitle),
		fmt.Sprintf("HTTP 状态码: %d", result.StatusCode),
		fmt.Sprintf("响应体大小: %d bytes", result.ContentLength),
	)

	// 风险判断
	if result.Risk {
		level = "ALERT"
		lines = append(lines, "风险判断: [x] 不允许公网访问，但实际可访问")
	} else {
		// 非风险状态
		if result.StatusCode != 200 {
			level = "INFO"
			lines = append(lines, "状态说明: [√] 非200状态码，可能未授权访问或资源受限")
		} else {
			level = "OK"
			lines = append(lines, "状态说明: 正常")
		}
	}

	return
}
