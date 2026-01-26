package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"os"
	"strings"
	"time"
)

type FeishuResponse struct {
	Msg string `json:"msg"`
}

func callAPI(webhook string, body []byte) error {

	client := req.C()

	// 为HTTP客户端设置一个“总超时时间“
	// 如果一次 HTTP 请求在 10 秒内没有完成（无论卡在哪个阶段），就直接失败并返回错误。
	client.SetTimeout(10 * time.Second)

	// 发送请求
	request := client.R()
	request.SetBody(body)
	request.SetHeader("Content-Type", "application/json")
	response, err := request.Post(webhook)

	log.Println("开始发送告警消息")

	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	// 获取响应状态码
	//log.Printf("状态码: %d\n", response.StatusCode)

	// 获取响应头
	//log.Printf("响应头: %+v\n", response.Header)

	// 获取响应体（自动转换为字符串）
	respBody := response.String()
	fmt.Printf("响应体: %s\n", respBody)

	var feishuResp FeishuResponse

	if err := json.Unmarshal([]byte(respBody), &feishuResp); err != nil {
		return fmt.Errorf("反序列失败：: %w", err)
	}
	log.Printf("飞书接口响应状态：%s", feishuResp.Msg)
	if feishuResp.Msg == "success" {
		log.Println("告警消息发送成功!")
	}

	return nil
}

// 从指定文件读取URL并返回字符串
func readURLFromFile(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 转换为字符串并去掉首尾空白符
	url := strings.TrimSpace(string(data))
	return url, nil
}

// 发送飞书扫描通知卡片：startMsg 可以是 "扫描开始" 或 "扫描结束"
func SendScanNotification(startMsg string) {
	webHookApi, err := readURLFromFile("config/feishuapi.txt")
	if err != nil {
		log.Fatalf("缺少飞书通知api配置文件：feishuapi.txt")
	}

	// 构造消息内容
	markdownContent := fmt.Sprintf("**<font color='blue'>🚀 %s</font>**", startMsg)

	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"schema": "2.0",
			"config": map[string]interface{}{
				"update_multi": true,
			},
			"body": map[string]interface{}{
				"direction": "vertical",
				"elements": []interface{}{
					map[string]interface{}{
						"tag":                "column_set",
						"flex_mode":          "stretch",
						"horizontal_spacing": "12px",
						"horizontal_align":   "left",
						"columns": []interface{}{
							map[string]interface{}{
								"tag":              "column",
								"width":            "weighted",
								"background_style": "bg-white",
								"elements": []interface{}{
									map[string]interface{}{
										"tag":        "markdown",
										"content":    markdownContent,
										"text_align": "left",
									},
								},
								"padding":            "12px 12px 12px 12px",
								"direction":          "vertical",
								"horizontal_spacing": "8px",
								"vertical_spacing":   "4px",
								"horizontal_align":   "left",
								"vertical_align":     "top",
								"weight":             1,
							},
						},
						"margin": "0px 0px 0px 0px",
					},
				},
			},
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": "🔍 安全扫描通知",
				},
				"subtitle": map[string]string{
					"tag":     "plain_text",
					"content": "",
				},
				"template": "blue",
				"padding":  "12px 8px 12px 8px",
			},
		},
	}

	// 序列化为 JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	callAPI(webHookApi, body)

}

func SendScanResult(url string, targetName string, title string, description string, time string) {
	// 飞书 Webhook 地址
	webHookApi, err := readURLFromFile("config/feishuapi.txt")
	if err != nil {
		log.Fatalf("config目录下缺少飞书通知api配置文件：feishuapi.txt")
	}

	assetName := targetName
	targetURL := url
	riskDesc := description
	checkTime := time
	assetTitle := title

	// 通知卡片大标题
	cardTitle := "公网访问告警"

	// 通知卡片副标题
	cardSubtitle := "受限服务可公网访问"

	// 通知卡片主体内容
	markdownContent := "**资产所属** ：" + assetName + "\n" +
		"**资产标题** ：" + assetTitle + "\n" +
		"**访问地址** ：" + targetURL + "\n" +
		"**检测时间** ：" + checkTime + "\n" +
		"**风险说明** ：" + riskDesc + "\n"

	//JSON 原始结构
	/*{
		"msg_type": "interactive",
		"card": {
		"schema": "2.0",
			"header": {...},
		"body": {...},
		"config": {...}
		}
	}*/

	// ===== 构造飞书卡片请求体=====
	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"schema": "2.0",
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": cardTitle,
				},
				"subtitle": map[string]string{
					"tag":     "plain_text",
					"content": cardSubtitle,
				},
				"template": "blue",
				"padding":  "12px 12px 12px 12px",
			},
			"body": map[string]interface{}{
				"direction": "vertical",
				"padding":   "12px 12px 12px 12px",
				"elements": []interface{}{
					map[string]interface{}{
						"tag":        "markdown",
						"content":    markdownContent,
						"text_align": "left",
						"text_size":  "normal_v2",
						"margin":     "0px 0px 0px 0px",
					},
				},
			},
			"config": map[string]interface{}{
				"update_multi": true,
				"style": map[string]interface{}{
					"text_size": map[string]interface{}{
						"normal_v2": map[string]string{
							"default": "normal",
							"pc":      "normal",
							"mobile":  "heading",
						},
					},
				},
			},
		},
	}

	// JSON 序列化
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// 发起请求
	callAPI(webHookApi, body)

}
