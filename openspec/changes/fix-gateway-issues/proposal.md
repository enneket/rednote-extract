## Why

当前 xhs-client-shell 的 HTTP Gateway 存在三个关键缺陷：(1) Spider_XHS 路径硬编码，production 部署必然断裂；(2) Cookie 管理双重标准，前端无法感知 cookies 实际内容；(3) API 响应结构不一致，用户笔记和搜索仅返回 URL，无法展示有效信息。

## What Changes

- **xhs-gateway-config**: Gateway 配置体系，支持 pip 包路径 / Docker 路径 / 本地开发路径多种部署模式
- **xhs-cookie-vault**: Cookie 统一管理，前端可读写完整 cookies 内容，不再依赖 Spider_XHS 目录
- **xhs-api-response**: API 响应结构对齐，`/user` 和 `/search` 也返回完整笔记信息（标题/点赞/收藏等）
- **xhs-data-export**: Excel 导出端点，接入 Spider_XHS 的 `save_to_xlsx` 实现下载

## Capabilities

### New Capabilities

- **xhs-gateway-config**: Gateway 启动时自动探测 Spider_XHS 路径，支持环境变量覆盖
- **xhs-cookie-vault**: Cookie 加密存储在 Gateway 可控目录，前端可查看已配置状态和内容摘要
- **xhs-api-response**: 笔记 API 统一返回结构化数据：`{id, title, author, liked, collected, commented, url, time}`
- **xhs-data-export**: `POST /api/v1/notes/export` 端点，接收笔记数据并生成 Excel 文件下载

### Modified Capabilities

- `xhs-http-gateway`: `/user` 和 `/search` 端点响应格式变更（从 URL 列表改为完整笔记信息）

## Impact

- `apps/gateway/app/core/spider_client.py`: 重构路径探测逻辑
- `apps/gateway/app/api/cookies.py`: Cookie 存储迁移，响应格式变更
- `apps/gateway/app/api/notes.py`: 统一响应结构，新增 export 端点
- `apps/client/src/api/client.ts`: 更新 TypeScript 类型以匹配新响应格式
