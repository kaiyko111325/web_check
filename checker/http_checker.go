package checker

import (
	"fmt"
	"github.com/imroc/req/v3"
	"regexp"
	"strings"
	"time"
)

// 跳转信息
type RedirectStep struct {
	URL        string
	Status     int
	Title      string
	Redirect   string
	Accessible bool // 是否可访问
}

// 跳转链检测器
type RedirectChecker struct {
	client *req.Client
}

// 创建新检测器
func NewRedirectChecker(timeout time.Duration) *RedirectChecker {
	client := req.C()
	client.SetTimeout(timeout)
	client.SetRedirectPolicy(req.NoRedirectPolicy())

	return &RedirectChecker{
		client: client,
	}
}

// 正则匹配 JS / Meta 重定向
var (
	// 完整匹配：location.replace('https://example.com/redirect')
	// 捕获组：https://example.com/redirect
	jsRedirect = regexp.MustCompile(`(?i)location\.replace\(['"]([^'"]+)['"]\)`)

	// 完整匹配：<meta http-equiv="refresh" content="0; url=https://example.com">
	// 捕获组：https://example.com
	metaRedirect = regexp.MustCompile(`(?i)<meta[^>]+http-equiv=["']refresh["'][^>]+content=["']\d+;\s*url=([^"']+)['"]`)

	// 完整匹配：<title>欢迎页面</title>
	// 捕获组：欢迎页面
	titleRegexp = regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)

	maxRedirects = 10
)

// 通过 HTTP 状态码是否可访问
func isAccessibleStatus(code int) bool {
	switch code {
	case 200:
		return true
	default:
		return false
	}
}

// 提取 title
func extractTitle(html string) string {
	matches := titleRegexp.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func isDeniedPage(body string) bool {
	// 是否明确是“访问被拒绝页面”
	isDeniedPage :=
		strings.Contains(body, "不允许访问") ||
			strings.Contains(body, "access denied") ||
			strings.Contains(body, "forbidden")
	return isDeniedPage
}

func isWAFBlockedResponse(body string) bool {
	return strings.Contains(body, "IP:") &&
		strings.Contains(body, "不允许访问")
}

func isBusinessResponse(body string) bool {
	return len(body) > 35 &&
		!isWAFBlockedResponse(body)
}

// / Check 检测目标 URL 的跳转链，并判断是否存在“未授权可访问”风险
func (r *RedirectChecker) Check(
	url string,
	targetName string,
	allowPublic bool,
) *ScanResult {

	result := &ScanResult{
		TargetName:  targetName,
		StartURL:    url,
		AllowPublic: allowPublic,
	}

	var chain []RedirectStep
	currentURL := url

	for i := 0; i < maxRedirects; i++ {

		resp, err := r.client.R().
			SetHeader("User-Agent", "Mozilla/5.0 (RedirectChecker)").
			Get(currentURL)

		if err != nil {
			// 请求失败
			chain = append(chain, RedirectStep{
				URL:        currentURL,
				Status:     0,
				Accessible: false,
			})
			result.Steps = chain
			result.StatusCode = 0
			result.Error = fmt.Errorf("请求失败：%s", TranslateError(err))
			result.Accessible = false
			return result
		}

		body := resp.String()
		title := extractTitle(body)

		step := RedirectStep{
			URL:    currentURL,
			Status: resp.StatusCode,
			Title:  title,
		}

		// 只有 200 才算可访问
		step.Accessible = isAccessibleStatus(resp.StatusCode)

		var nextURL string

		// HTTP 3xx 跳转
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			nextURL = resp.Header.Get("Location")
			step.Redirect = nextURL
		}

		// JS 跳转
		if nextURL == "" {
			if jsMatch := jsRedirect.FindStringSubmatch(body); len(jsMatch) > 1 {
				nextURL = jsMatch[1]
				step.Redirect = nextURL
			}
		}

		// Meta Refresh 跳转
		if nextURL == "" {
			if metaMatch := metaRedirect.FindStringSubmatch(body); len(metaMatch) > 1 {
				nextURL = metaMatch[1]
				step.Redirect = nextURL
			}
		}

		chain = append(chain, step)

		if nextURL == "" {
			// 扫描结束
			result.Steps = chain
			result.FinalURL = currentURL
			result.FinalTitle = title
			result.StatusCode = resp.StatusCode

			result.Accessible = isAccessibleStatus(resp.StatusCode)

			result.ContentLength = int64(len(resp.Bytes()))

			if isWAFBlockedResponse(body) {
				// 被 WAF 正常拦截，不属于风险
				result.Risk = false
				return result
			}

			// 真正的业务响应 + 不允许公网
			if isBusinessResponse(body) && result.Accessible && !result.AllowPublic && !isDeniedPage(body) {
				result.Risk = true
				return result
			}

			return result
		}

		currentURL = nextURL
	}

	// 超过最大跳转
	result.Steps = chain
	result.Error = fmt.Errorf("跳转次数超过 %d，可能存在循环跳转", maxRedirects)
	return result
}
