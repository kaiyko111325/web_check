package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"web-exposure-check/FeiShuApi"
	"web-exposure-check/checker"
	"web-exposure-check/config"
	"web-exposure-check/logger"
	"web-exposure-check/notifier"
)

// 从飞书指定文档中读取待扫描的url，并保存为targets.yaml到config目录下
func GetFeiShuAuthCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization Failed",
		})
		return
	}

	// 同步生成配置文件
	if err := FeiShuApi.FeiShuSheeDateToYaml(code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 文件生成完成后，再异步扫描
	go AccessScan()

	c.JSON(http.StatusOK, gin.H{
		"success": "授权成功，扫描任务已启动",
	})
}

// ToLoginFeiShu 处理 /login 路径，访问该路径会跳转到飞书授权页面
func ToLoginFeiShu(c *gin.Context) {
	// 飞书 OAuth2 授权 URL
	authURL := "https://accounts.feishu.cn/open-apis/authen/v1/authorize" +
		"?client_id=cli_a9faf1d347385bcb" +
		"&response_type=code" +
		"&redirect_uri=http://113.44.78.107:18989/authCode" +
		"&scope=wiki:node:read wiki:wiki wiki:wiki:readonly sheets:spreadsheet drive:drive sheets:spreadsheet:readonly drive:drive:readonly sheets:spreadsheet:read" +
		"&state=RANDOMSTRING"

	// 重定向到飞书授权页面
	c.Redirect(http.StatusFound, authURL)
}

// 获取服务器公网ip，用于打印飞书授权地址
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me/ip") // 返回纯文本 IP
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 探测网站是否可访问的所有过程
func AccessScan() {
	logFile, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logFile.Close()

	cfg, err := config.LoadTargetsFromFile("config/targets.yaml")
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
		// 打印错误并退出
	}

	rc := checker.NewRedirectChecker(10 * time.Second)

	// 记录扫描开始时间
	startTime := time.Now().Format("2006-01-02 15:04:05")
	startMsg := "扫描开始：" + startTime + "\n"

	notifier.SendScanNotification(startMsg)

	for _, t := range cfg.Targets {

		result := rc.Check(t.URL, t.Name, t.AllowPublic)

		level, title, lines := checker.PrintResult(result)

		// 日志打印，同时如果告警level是 ALERT 就触发飞书告警
		logger.Print(result, logger.Level(level), title, lines...)
	}

	// 记录扫描结束时间
	endTime := time.Now().Format("2006-01-02 15:04:05")
	endMsg := "扫描结束：" + endTime + "\n"
	notifier.SendScanNotification(endMsg)

}

// -------------------- main --------------------
func main() {

	serverIP, err := GetPublicIP()
	if err != nil {
		fmt.Println("获取外网IP失败:", err)
		return
	}
	fmt.Println("本机外网 IP:", serverIP)

	// 创建 Gin 路由引擎
	r := gin.Default()

	// 飞书授权跳转
	r.GET("/login", ToLoginFeiShu)

	// 接收飞书回调函数的传参code
	r.GET("/authCode", GetFeiShuAuthCode)

	fmt.Printf("\033[31m⚠️  请访问下方 URL 完成飞书授权  ⚠️\033[0m\n")
	fmt.Printf("\033[31m👉 请访问：http://%s:18989/login 完成飞书授权\033[0m\n", serverIP)
	fmt.Printf("\033[31m⚠️  请访问上方 URL 完成飞书授权  ⚠️\033[0m\n")

	// 启动服务
	r.Run(":18989")

}
