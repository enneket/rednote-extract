## Why

Spider_XHS 是命令行驱动的 Python 爬虫库，核心用户是开发者通过代码调用。对于非技术背景的运营人员而言，直接使用 Python 代码门槛太高。将 Spider_XHS 封装为 Windows 原生客户端（exe），让运营人员可以图形化操作小红书数据采集。

## What Changes

- 新增 **xhs-desktop-client**：Windows 桌面客户端，提供图形化界面操作爬虫功能，打包为单文件 exe
- 新增 **xhs-http-gateway**：FastAPI HTTP 网关，封装 Spider_XHS 为 REST API，客户端通过 HTTP 调用

## Capabilities

### New Capabilities

- **xhs-desktop-client**: Windows exe 客户端（Tauri + React 或 PyInstaller + Flet），提供 Cookie 管理、笔记采集、用户采集、搜索采集的可视化界面
- **xhs-http-gateway**: FastAPI 服务，封装 `Data_Spider` 类为 RESTful API，客户端与 Spider_XHS 之间通过 HTTP 通信

### Modified Capabilities

- 无

## Impact

- 新增目录：`apps/client/`（桌面客户端）、`apps/gateway/`（HTTP 网关）
- 客户端打包为 Windows exe 分发，可独立运行
- HTTP Gateway 可同时供客户端和 AI Agent 调用
