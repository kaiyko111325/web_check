package FeiShuApi

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"strings"
	"time"
)

/************** 响应结构体 **************/

type SpreadsheetSheetsResponse struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data SpreadsheetSheetsData `json:"data"`
}

type SpreadsheetSheetsData struct {
	Sheets []SpreadsheetSheet `json:"sheets"`
}

type SpreadsheetSheet struct {
	SheetID      string `json:"sheet_id"`
	Title        string `json:"title"`
	Index        int    `json:"index"`
	Hidden       bool   `json:"hidden"`
	ResourceType string `json:"resource_type"`

	GridProperties GridProperties `json:"grid_properties"`
}

type GridProperties struct {
	RowCount          int `json:"row_count"`
	ColumnCount       int `json:"column_count"`
	FrozenRowCount    int `json:"frozen_row_count"`
	FrozenColumnCount int `json:"frozen_column_count"`
}

// GetSpreadsheetInfo
// 根据 spreadsheet_token 获取该表下所有工作表（sheet）信息
func GetSpreadsheetInfo(
	accessToken string,
	spreadsheetToken string,
) ([]SpreadsheetSheet, error) {

	if accessToken == "" || spreadsheetToken == "" {
		return nil, fmt.Errorf("accessToken or spreadsheetToken is empty")
	}

	url := fmt.Sprintf(
		"https://open.feishu.cn/open-apis/sheets/v3/spreadsheets/%s/sheets/query",
		spreadsheetToken,
	)

	client := req.C().
		SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+strings.TrimSpace(accessToken)).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"http request failed, status=%d, body=%s",
			resp.StatusCode,
			resp.String(),
		)
	}

	var result SpreadsheetSheetsResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("feishu api error: %s", result.Msg)
	}

	return result.Data.Sheets, nil
}

func GetSheetByTitle(
	sheets []SpreadsheetSheet,
	title string,
) (*SpreadsheetSheet, error) {

	if len(sheets) == 0 {
		return nil, fmt.Errorf("sheet list is empty")
	}

	for _, sheet := range sheets {
		if sheet.Title == title {
			return &sheet, nil
		}
	}

	return nil, fmt.Errorf("sheet with title %q not found", title)
}

func GetSheetInfoByID(
	sheets []SpreadsheetSheet,
	sheetID string,
) (*SpreadsheetSheet, error) {

	if sheetID == "" {
		return nil, fmt.Errorf("sheetID is empty")
	}

	if len(sheets) == 0 {
		return nil, fmt.Errorf("sheet list is empty")
	}

	for i := range sheets {
		if sheets[i].SheetID == sheetID {
			return &sheets[i], nil
		}
	}

	return nil, fmt.Errorf("sheet with id %q not found", sheetID)
}

func PrintSheetsInfo(sheets []SpreadsheetSheet) {

	if len(sheets) == 0 {
		fmt.Println("sheetsData is empty")
		return
	}

	fmt.Printf("Total Sheets: %d\n", len(sheets))
	fmt.Println("====================================")

	for i, sheet := range sheets {

		fmt.Printf("Sheet #%d\n", i)
		fmt.Printf("  SheetID       : %s\n", sheet.SheetID)
		fmt.Printf("  Title         : %s\n", sheet.Title)
		fmt.Printf("  Index         : %d\n", sheet.Index)
		fmt.Printf("  Hidden        : %v\n", sheet.Hidden)
		fmt.Printf("  ResourceType  : %s\n", sheet.ResourceType)

		fmt.Println("  GridProperties:")
		fmt.Printf("    RowCount          : %d\n", sheet.GridProperties.RowCount)
		fmt.Printf("    ColumnCount       : %d\n", sheet.GridProperties.ColumnCount)
		fmt.Printf("    FrozenRowCount    : %d\n", sheet.GridProperties.FrozenRowCount)
		fmt.Printf("    FrozenColumnCount : %d\n", sheet.GridProperties.FrozenColumnCount)

		fmt.Println("------------------------------------")
	}
}
