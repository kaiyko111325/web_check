package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"testfeishu/checker"
	"testfeishu/notifier"
	"time"
)

type Level string

const (
	LevelOK    Level = "OK"
	LevelALERT Level = "ALERT"
	LevelERROR Level = "ERROR"
)

// Print 打印日志，并在 ALER 时触发飞书告警
func Print(result *checker.ScanResult, level Level, title string, lines ...string) {
	prefix := fmt.Sprintf("[%s] %s", level, title)
	log.Println(prefix)

	for _, line := range lines {
		log.Println("  " + line)
	}

	// 如果是 ALERT 或 ERROR，则调用飞书
	if level == LevelALERT {

		description := fmt.Sprintf(
			"该资产可被外部访问，存在安全风险，响应体大小: %d bytes",
			result.ContentLength,
		)

		// 当前时间
		now := time.Now().Format("2006-01-02 15:04:05")

		// 调用飞书告警
		notifier.SendScanResult(result.StartURL, result.TargetName, result.FinalTitle, description, now)
	}
}

func InitLogger() (*os.File, error) {
	// 创建 logs 目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf(
		"logs/scan-%s.log",
		time.Now().Format("2006-01-02_15-04-05"),
	)

	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	// 同时输出到 控制台 + 文件
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	// 可选：统一日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return file, nil
}
