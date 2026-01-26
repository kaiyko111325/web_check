package FeiShuApi

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"strings"
	"time"
)

type SpreadsheetValuesResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SpreadsheetToken string             `json:"spreadsheetToken"`
		Revision         int                `json:"revision"`
		TotalCells       int                `json:"totalCells"`
		ValueRanges      []SpreadsheetRange `json:"valueRanges"`
	} `json:"data"`
}

type SpreadsheetRange struct {
	Range  string     `json:"range"`
	Values [][]string `json:"values"`
}

func GetSheetDataValues(
	accessToken string,
	spreadsheetToken string,
	sheetId string,
	rangeStr string,
) ([][]string, error) {

	if accessToken == "" || spreadsheetToken == "" || sheetId == "" || rangeStr == "" {
		return nil, fmt.Errorf("invalid parameter")
	}

	url := fmt.Sprintf(
		"https://open.feishu.cn/open-apis/sheets/v2/spreadsheets/%s/values_batch_get",
		spreadsheetToken,
	)

	// 拼 ranges 参数
	ranges := fmt.Sprintf("%s!%s", sheetId, rangeStr)

	client := req.C().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+strings.TrimSpace(accessToken)).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetQueryParam("ranges", ranges).
		SetQueryParam("valueRenderOption", "ToString").
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.String())
	}

	var result SpreadsheetValuesResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("feishu api error: %s", result.Msg)
	}

	if len(result.Data.ValueRanges) == 0 {
		return nil, fmt.Errorf("no valueRanges returned")
	}

	return result.Data.ValueRanges[0].Values, nil
}
