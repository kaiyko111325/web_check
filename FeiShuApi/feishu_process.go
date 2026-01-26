package FeiShuApi

import (
	"fmt"
	"os"
	"path"
	"web-exposure-check/SheetDataToYaml"
)

// -------------------- 步骤函数 --------------------

// 获取 AccessToken
func GetAccessToken(code string) (string, error) {
	tokenResp, err := GetAccessTokenFromFeiShu(code)
	if err != nil {
		return "", fmt.Errorf("access token fail: %w", err)
	}
	return tokenResp.AccessToken, nil
}

// 获取 spreadsheetToken
func GetSpreadsheetToken(accessToken, wikiURL string) (string, error) {
	nodeToken := path.Base(wikiURL)
	fmt.Printf("nodeToken: %s\n", nodeToken)

	knowledgeNode, err := GetKnowledgeSpaceNodeInfo(accessToken, nodeToken)
	if err != nil {
		return "", err
	}

	fmt.Printf("ObjToken（spreadsheetToken）: %s\n", knowledgeNode.Node.ObjToken)
	return knowledgeNode.Node.ObjToken, nil
}

// 获取 Sheets 信息
func GetSheets(accessToken, spreadsheetToken string) ([]SpreadsheetSheet, error) {
	sheetsData, err := GetSpreadsheetInfo(accessToken, spreadsheetToken)
	if err != nil {
		return nil, err
	}

	PrintSheetsInfo(sheetsData)
	return sheetsData, nil
}

// 根据 sheet title 获取 SheetID
func GetSheetIDByTitle(sheets []SpreadsheetSheet, title string) (string, error) {
	sheetData, err := GetSheetByTitle(sheets, title)
	if err != nil {
		return "", err
	}
	return sheetData.SheetID, nil
}

// 获取 sheet 数据
func GetSheetValues(accessToken, spreadsheetToken, sheetID, sheetRange string) ([][]string, error) {
	values, err := GetSheetDataValues(accessToken, spreadsheetToken, sheetID, sheetRange)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// 保存为 YAML
func SaveValuesToYAML(values [][]string, filename string) error {
	return SheetDataToYaml.SaveSheetValuesToYAML(values, filename)
}

func FeiShuSheeDateToYaml(code string) error {

	wikiURL := "https://youzu.feishu.cn/wiki/AuDWwVPh7irIfIkkv8ZcK7Hlnth"
	sheetTitle := "汇总数据（公网GM）"
	sheetRange := "F:G"
	yamlFile := "config/targets.yaml" // ✅ 相对路径

	// 确保目录存在
	if err := os.MkdirAll("config", 0755); err != nil {
		return fmt.Errorf("创建 config 目录失败: %w", err)
	}

	// 1. 获取 AccessToken
	accessToken, err := GetAccessToken(code)
	if err != nil {
		return fmt.Errorf("获取 AccessToken 失败: %w", err)
	}

	// 2. 获取 spreadsheetToken
	spreadsheetToken, err := GetSpreadsheetToken(accessToken, wikiURL)
	if err != nil {
		return fmt.Errorf("获取 spreadsheetToken 失败: %w", err)
	}

	// 3. 获取 sheets 信息
	sheetsData, err := GetSheets(accessToken, spreadsheetToken)
	if err != nil {
		return fmt.Errorf("获取 sheets 信息失败: %w", err)
	}

	// 4. 获取 sheetID
	sheetID, err := GetSheetIDByTitle(sheetsData, sheetTitle)
	if err != nil {
		return fmt.Errorf("获取 SheetID 失败: %w", err)
	}
	fmt.Printf("%s 的 SheetID：%s\n", sheetTitle, sheetID)

	// 5. 获取 sheet 数据
	values, err := GetSheetValues(accessToken, spreadsheetToken, sheetID, sheetRange)
	if err != nil {
		return fmt.Errorf("获取 SheetValues 失败: %w", err)
	}

	// 6. 打印数据（可选）
	for i, row := range values {
		if len(row) >= 2 {
			fmt.Printf("Row %d: URL=%s | Title=%s\n", i, row[0], row[1])
		}
	}

	// 7. 保存为 YAML
	if err := SaveValuesToYAML(values, yamlFile); err != nil {
		return fmt.Errorf("保存 YAML 失败: %w", err)
	}

	fmt.Println("YAML 文件保存完成:", yamlFile)
	return nil
}
