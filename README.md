# RedNote Extract (小红书数据提取与二创工具)

RedNote Extract 是一个基于 Go 语言开发的小红书数据抓取与 AI 二创工具。它能够自动化地根据关键词抓取小红书笔记，并利用大语言模型（LLM）对内容进行分析和重新创作，生成高质量的新笔记。

## ✨ 主要功能

- **🔍 关键词搜索抓取**: 支持自定义关键词，批量抓取小红书笔记数据（标题、正文、评论等）。
- **🤖 AI 智能二创**: 集成 LLM（支持 Qwen, OpenAI 等），对抓取的素材进行深度理解和二次创作。
- **🌐 浏览器自动化**: 基于 Playwright，模拟真实用户行为，支持 Cookie 注入以获取更高权限数据。
- **📂 本地数据存储**: 自动将抓取的原始数据和生成的新笔记保存为 JSON 文件，便于后续分析。

## 🛠️ 技术栈

- **语言**: Go 1.24+
- **浏览器自动化**: [Playwright Go](https://github.com/playwright-community/playwright-go)
- **AI 框架**: [CloudWeGo Eino](https://github.com/cloudwego/eino)
- **HTTP 客户端**: Resty
- **配置管理**: Viper

## 🚀 快速开始

### 1. 环境准备

确保你的系统已安装 Go 1.24 或更高版本。

### 2. 克隆项目

```bash
git clone https://github.com/enneket/rednote-extract.git
cd rednote-extract
```

### 3. 安装依赖

下载 Go 模块依赖：

```bash
go mod download
```

安装 Playwright 浏览器及其依赖（首次运行必须执行）：

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
```

### 4. 配置文件

复制示例配置文件并进行修改：

```bash
cp env.example .env
```

编辑 `.env` 文件，填入你的配置信息：

- **LLM 配置**: 设置 `LLM_API_KEY` 和 `LLM_API_BASE_URL`（例如使用通义千问 Qwen）。
- **Cookie**: 如果需要登录状态抓取，请在 `COOKIES` 字段填入你的小红书 Cookie。
- **其他**: 可调整 `MAX_NOTES`（抓取数量）等参数。

### 5. 运行项目

使用 `-keywords` 参数指定要抓取的关键词：

```bash
go run cmd/main.go -keywords "编程副业"
```

## 📂 输出说明

运行后，程序会在 `data/xhs/` 目录下生成以“时间戳_关键词”命名的文件夹，例如 `data/xhs/20260121_150405_编程副业/`。

文件夹内包含：
- 抓取的原始笔记数据（JSON 格式）。
- AI 生成的新笔记内容。

## ⚠️ 免责声明

本项目仅供学习和研究使用。请勿用于任何商业用途或违反小红书平台规则的行为。使用本工具产生的任何后果由使用者自行承担。
