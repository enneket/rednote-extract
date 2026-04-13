## Context

Spider_XHS 是一个纯 Python 爬虫库，通过 `Data_Spider` 类提供小红书数据采集能力。当前只能通过 Python 代码调用，运营人员无法直接使用。需要封装为 Windows exe 客户端，降低使用门槛。

## Goals / Non-Goals

**Goals:**
- 运营人员可通过图形界面操作 Spider_XHS（配置 Cookie、触发采集、查看结果）
- 客户端打包为单文件 Windows exe，分发给非技术用户
- 保留 HTTP API 能力供 AI Agent 调用
- 客户端与爬虫逻辑解耦，便于独立维护

**Non-Goals:**
- 不修改 Spider_XHS 源码
- 不实现小红书内容发布功能（仅数据采集）
- 不做移动端适配（桌面端全平台覆盖即可）

## Decisions

### 1. 客户端技术栈：Tauri + React + TypeScript

**决定：** 使用 Tauri + React 构建 UI，跨平台打包（Windows/macOS/Linux）。

**理由：**
- Tauri 原生支持跨平台，打包体积小（<20MB exe），启动快
- React + TypeScript 提供现代 Web 化 UI
- 与 Spider_XHS Python 生态通过 HTTP Gateway 解耦
- 可同时输出 Windows exe、macOS dmg、Linux AppImage

**替代方案：**
- Electron + Python：打包体积过大（>200MB），内存占用高
- PyInstaller + Flet：仅 Windows，单文件体积大（100-200MB）
- Flutter Desktop：需要 Dart 技能栈，与现有 Python 生态割裂

### 2. 客户端与爬虫解耦：HTTP Gateway 模式

**决定：** 客户端不直接调用 Spider_XHS，而是通过 FastAPI HTTP Gateway 间接调用。

**理由：**
- 客户端进程与爬虫进程隔离，爬虫崩溃不影响 UI
- Gateway 可独立部署，同时供客户端和 AI Agent 使用
- 便于未来扩展更多客户端（CLI、Web、移动端）

**替代方案：**
- 直接在客户端内调用 Spider_XHS：简单但无法复用，爬虫逻辑与 UI 耦合

### 3. Gateway 与 Spider_XHS 集成：Python 子进程

**决定：** Gateway 通过 `subprocess` 调用 Spider_XHS 的 Python 脚本，进程间通过 JSON 文件或 HTTP API 通信。

**理由：**
- 不需要修改 Spider_XHS 源码
- Spider_XHS 有自己的 Cookie 管理和初始化逻辑，保持独立运行

### 4. 客户端打包架构

```
┌─────────────────────────────────────────────┐
│  跨平台桌面客户端 (Tauri 打包)                │
│  ┌─────────────────┐   ┌────────────────┐  │
│  │  React + TS UI  │──▶│  HTTP Gateway   │  │
│  │  (Web 前端)      │   │  (FastAPI)      │  │
│  └─────────────────┘   └────────┬───────┘  │
│                                │           │
│                       ┌────────▼────────┐   │
│                       │  Spider_XHS    │   │
│                       │  (Python 子进程) │   │
│                       └────────────────┘   │
└─────────────────────────────────────────────┘
          │              │              │
       Windows         macOS          Linux
       .exe           .dmg          .AppImage
```

All-in-one exe 内部：Gateway 和 Spider_XHS 作为嵌入资源运行。

## Risks / Trade-offs

- **打包体积**：PyInstaller 单文件约 100-200MB，比 Tauri 大。→ 接受，因为开发效率优先
- **启动速度**：首次启动需要解压 Python 环境，约 5-10 秒 → 可通过引导界面改善体验
- **Cookie 安全**：Cookie 存储在本地文件，加密存储 → 实现中需引入加密方案
