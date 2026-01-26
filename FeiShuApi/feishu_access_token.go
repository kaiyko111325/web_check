package FeiShuApi

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"strings"
)

const (
	GetAccessTokenURL = "https://open.feishu.cn/open-apis/authen/v2/oauth/token"
	ClientID          = "cli_a9faf1d347385bcb"
	ClientSecret      = "9VXZ6Q70l0KG4KPVWueeYhlONSNX116n"
	RedirectURI       = "http://113.44.78.107:18989/authCode"
)

type FeishuAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

type FeishuAccessTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Code        int    `json:"code"`
}

/*
GetAccessTokenFromFeiShu
使用授权 code 换取飞书 access_token
*/
func GetAccessTokenFromFeiShu(code string) (*FeishuAccessTokenResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("authorization code is empty")
	}

	reqBody := FeishuAccessTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Code:         code,
		RedirectURI:  RedirectURI,
	}

	// 发送请求
	client := req.C()
	request := client.R()
	request.SetBody(reqBody)
	request.SetHeader("Content-Type", "application/json; charset=utf-8")

	resp, err := request.Post(GetAccessTokenURL)
	if err != nil {
		return nil, err
	}

	//fmt.Println(resp.String())

	// 5. 判断 HTTP 状态码是否正常
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"http request failed, status=%d, body=%s",
			resp.StatusCode,
			resp.String(),
		)
	}

	// 6. 手动解析 JSON 响应
	var tokenResp FeishuAccessTokenResponse
	if err := json.Unmarshal(resp.Bytes(), &tokenResp); err != nil {
		return nil, err
	}

	fmt.Printf("tokenResp.AccessToken: %s\n", tokenResp.AccessToken)
	// 7. 判断飞书返回的业务状态
	if tokenResp.Code != 0 {
		return nil, fmt.Errorf("feishu oauth failed")
	}

	// 8. 返回解析结果
	return &tokenResp, nil
}
