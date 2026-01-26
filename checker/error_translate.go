package checker

import (
	"errors"
	"golang.org/x/net/context"
	"net"
	"strings"
)

// TranslateError 将 Go 原始错误翻译为用户可理解的失败原因
func TranslateError(err error) string {
	if err == nil {
		return ""
	}

	// 1. context 超时
	if errors.Is(err, context.DeadlineExceeded) {
		return "请求超时（目标长时间未响应，可能端口被防火墙拦截或服务不可达）"
	}

	// 2. net.Error（包括各种 timeout）
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return "网络请求超时（目标未在规定时间内响应）"
		}
		return "网络错误（连接失败或目标不可达）"
	}

	// 3. 根据错误文本判断
	msg := err.Error()

	switch {
	case strings.Contains(msg, "Client.Timeout"):
		return "请求超时（服务器响应过慢）"
	case strings.Contains(msg, "connection refused"):
		return "连接被拒绝（端口未开放或服务未运行）"
	case strings.Contains(msg, "no such host"):
		return "域名解析失败（DNS 无法解析目标）"
	case strings.Contains(msg, "i/o timeout"):
		return "网络 I/O 超时（可能被防火墙拦截）"
	}

	// 4. 最终还是没有匹配到错误，给一个默认报错
	return "请求失败（网络异常或目标不可达）"
}
