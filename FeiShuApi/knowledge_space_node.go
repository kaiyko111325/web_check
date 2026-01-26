package FeiShuApi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"strings"
	"time"
)

// 响应体结构
// WikiNode 单独定义，复用性高
type WikiNode struct {
	Creator         string `json:"creator"`
	HasChild        bool   `json:"has_child"`
	NodeCreateTime  string `json:"node_create_time"`
	NodeCreator     string `json:"node_creator"`
	NodeToken       string `json:"node_token"`
	NodeType        string `json:"node_type"`
	ObjCreateTime   string `json:"obj_create_time"`
	ObjEditTime     string `json:"obj_edit_time"`
	ObjToken        string `json:"obj_token"`
	ObjType         string `json:"obj_type"`
	OriginNodeToken string `json:"origin_node_token"`
	OriginSpaceID   string `json:"origin_space_id"`
	Owner           string `json:"owner"`
	ParentNodeToken string `json:"parent_node_token"`
	SpaceID         string `json:"space_id"`
	Title           string `json:"title"`
}

// WikiNodeData 包含 Node
type WikiNodeData struct {
	Node WikiNode `json:"node"`
}

// WikiNodeResponse 飞书返回的完整响应结构
type WikiNodeResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data WikiNodeData `json:"data"`
}

// GetKnowledgeSpaceNodeInfo 调用飞书 Wiki API 获取知识空间节点信息
func GetKnowledgeSpaceNodeInfo(accessToken string, nodeToken string) (*WikiNodeData, error) {
	if accessToken == "" || nodeToken == "" {
		return nil, fmt.Errorf("access_token or node_token is empty")
	}

	// 1. 构造请求 URL
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/wiki/v2/spaces/get_node?token=%s&obj_type=wiki", nodeToken)

	// 2. 创建 req 客户端
	client := req.C()
	client.SetTimeout(10 * time.Second)
	// 设置 HTTP/HTTPS 代理
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request := client.R()

	request.SetHeader("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	request.SetHeader("Content-Type", "application/json; charset=utf-8")

	// 打印请求 URL
	fmt.Println("Request URL:", url)

	// 打印请求 Header
	fmt.Println("Request Headers:")
	for k, v := range request.Headers {
		fmt.Printf("%s: %v\n", k, v)
	}

	// 3. 发送 GET 请求
	resp, err := request.Get(url)

	if err != nil {
		return nil, err
	}

	fmt.Printf("response code %d\n", resp.StatusCode)

	// 4. 判断 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed, status=%d, body=%s",
			resp.StatusCode, resp.String())
	}

	// 5. 手动解析 JSON
	var wikiResp WikiNodeResponse
	if err := json.Unmarshal(resp.Bytes(), &wikiResp); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 6. 判断飞书返回的业务状态
	if wikiResp.Code != 0 {
		return nil, fmt.Errorf("feishu wiki api failed, msg=%s", wikiResp.Msg)
	}

	// 7. 返回节点数据
	return &wikiResp.Data, nil
}
