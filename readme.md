# Web Exposure Checker（飞书资产联动版）

一个基于 **Go** 开发的 Web 资产公网可访问性检测系统。  
系统通过 **飞书 OAuth 登录** 获取用户授权，从 **飞书电子表格**中自动拉取资产信息，生成扫描目标配置，并对资产进行可访问性探测，在发现违规公网暴露时通过 **飞书机器人**进行告警。

---

## 一、整体流程说明

系统整体工作流程如下：

1. 启动程序
2. 用户访问指定路径  
3. 跳转飞书进行 OAuth 授权
4. 授权成功后：

- 使用 Access Token 调用飞书开放接口
- 从 **指定飞书电子表格**中读取资产数据

5. 自动生成 `config/targets.yaml`
6. 执行 Web 资产公网可访问性扫描
7. 若发现风险资产：

- 通过飞书机器人发送告警消息

---

## 二、核心功能

### 1️⃣ 飞书 OAuth 登录

- 使用飞书 **OAuth 2.0** 授权机制
- 通过 `/login` 路由触发认证流程
- 授权成功后自动进入资产同步和扫描流程

### 2️⃣ 飞书电子表格资产同步

- 从飞书「多维表格 / 电子表格」中读取资产信息
- 支持字段示例：
- 资产名称
- 访问地址（URL）
- 是否允许公网访问
- 资产负责人（可选）
- 自动转换为本地扫描配置文件

### 3️⃣ 自动生成 targets.yaml

授权成功后自动生成：

```yaml
targets:
- name: Fancy Data
 url: http://fancy-gm.uuzuonline.net
 allow_public: false
```

无需人工维护扫描目标。

### 4️⃣ Web 公网可访问性探测

- 支持 HTTP / HTTPS
- 支持跳转链跟踪（3xx）
- 仅 `HTTP 200` 视为“可访问”
- 自动识别并排除：
  - WAF 拦截响应（如：`IP:xxx 不允许访问`）
  - 非业务错误页面

### 5️⃣ 风险识别与告警

当满足以下条件时触发告警：

- 实际可访问（HTTP 200）
- 资产标记为 **不允许公网访问**
- 响应内容为真实业务页面

告警内容通过 **飞书机器人**发送。

------

## 三、项目结构

```
web-exposure-checker/
│  go.mod
│  go.sum
│  main.go               			# 程序入口
│  readme.md
│
├─checker							# 核心扫描与风险判断逻辑
│      error_translate.go
│      http_checker.go
│      printer.go
│      result.go
│
├─config
│      feishuapi.txt
│      loader.go
├─ log                   			# 扫描日志
├─FeiShuApi							# 飞书读取表格相关逻辑
│      feishu_access_token.go
│      feishu_process.go
│      knowledge_space_node.go
│      sheet_data.go
│      sheet_info.go
│      work_sheet.go
│
├─logger
│      logger.go					# 记录探测日志
│
├─logs
├─notifier
│      feishu.go					# 飞书机器人 Webhook
│
├─picture
│      PixPin_2026-01-26_11-24-26.png
│      PixPin_2026-01-26_11-26-26.png
│
└─SheetDataToYaml					# 飞书表格数据转化为target.yaml格式
        convert_sheet_data_to_yaml.go
```

------

## 四、配置说明

### 1️⃣ 飞书机器人配置

```
config/feishuapi.txt
https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxx-xxxx
```

用于发送风险告警消息。

------

### 2️⃣ 飞书电子表格要求

电子表格需至少包含以下字段（列名可配置）：

| 字段名 | 说明     |
| ------ | -------- |
| name   | 资产名称 |
| url    | 访问地址 |

### 1️⃣ 编译或直接运行

```bash
go mod tidy
go run main.go
```

### 2️⃣ 访问登录地址

浏览器打开：

```
http://<公网IP>:18989/login
```



![PixPin_2026-01-26_11-24-26](D:\go_project\web-exposure-check\picture\PixPin_2026-01-26_11-24-26.png)





![PixPin_2026-01-26_11-26-26](D:\go_project\web-exposure-check\picture\PixPin_2026-01-26_11-26-26.png)

### 3️⃣ 飞书授权成功后

系统将自动完成：

- 资产拉取
- `targets.yaml` 生成
- 扫描执行
- 飞书告警（如存在风险）

------

## 六、告警示例（飞书）

```
【公网暴露风险告警】

资产标题：baidu
访问地址：http://baidu.com
检测时间：2026-01-19 14:17:55
风险说明：该资产可被外部访问，存在安全风险
响应体大小：5123 bytes
```

------

## 七、设计特点

- **自动化**：无需手工维护扫描资产
- **强可扩展**：资产来源、告警方式可扩展
- **真实业务识别**：避免 WAF / 错误页误判
- **工程化结构**：适合企业内部持续运行

------

## 八、注意事项

- 请确保对扫描资产拥有合法授权
- 建议部署在内网服务器并配置公网访问策略
- 飞书 Access Token 需注意有效期和权限范围

------

