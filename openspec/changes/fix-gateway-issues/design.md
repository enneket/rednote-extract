## Context

xhs-client-shell 的 HTTP Gateway 已完成核心功能，但存在三个 production-ready 隐患需要修复。

## Goals / Non-Goals

**Goals:**
- Spider_XHS 路径探测：从硬编码绝对路径改为多模式自动探测 + 环境变量覆盖
- Cookie 统一管理：迁移到 Gateway 可控目录，前端可见实际内容
- API 响应对齐：用户笔记和搜索端点返回完整笔记数据而非 URL 列表
- Excel 导出：实现笔记数据下载

**Non-Goals:**
- 不改变现有 API 端点路径（保持兼容）
- 不修改 Spider_XHS 源码

## Decisions

### 1. Spider_XHS 路径探测：优先级 + 环境变量

**决定：** 按优先级探测 Spider_XHS 路径，支持环境变量覆盖。

**探测顺序：**
1. 环境变量 `SPIDER_XHS_PATH`（最高优先级，Docker / production 使用）
2. pip 包安装路径 `importlib.util.find_spec("spider_xhs")`
3. Gateway 同级目录 `../Spider_XHS`（开发模式）
4. 抛出明确错误，提示用户配置

**理由：** 环境变量覆盖是最灵活的 production 方案，pip 包安装是标准化分发路径，同级目录兼容开发体验。

### 2. Cookie 存储：迁移到 Gateway 可控目录

**决定：** Cookie 存储从 `Spider_XHS/.cookies` 迁移到 `~/.xhs-gateway/cookies.enc`。

**结构：**
```
~/.xhs-gateway/
  cookies.enc   # AES-256-GCM 加密存储
  config.yaml   # Gateway 配置（Spider_XHS 路径等）
```

**理由：** 桌面客户端独立运行时不应依赖 Spider_XHS 目录。加密存储保障 Cookie 安全。

**备选方案：** 直接存 `cookies.json` 明文 → 拒绝，安全风险。

### 3. API 统一响应结构

**决定：** 所有笔记 API 返回统一的笔记信息结构。

**统一响应字段：**
```json
{
  "id": "note_id",
  "title": "标题",
  "author": "作者昵称",
  "author_id": "user_id",
  "type": "图集|视频",
  "liked": 1234,
  "collected": 567,
  "commented": 89,
  "shared": 12,
  "url": "https://...",
  "tags": ["标签1", "标签2"],
  "time": "2024-01-01 12:00:00",
  "ip_location": "广东"
}
```

**`/single` 端点：** 直接复用 `handle_note_info()` 返回值，取其中上述字段映射。

**`/user` 和 `/search` 端点：** 遍历 URL 列表，逐条调用 `get_note_info` 获取完整数据，返回统一结构数组。

**性能考量：** 批量采集时串行调用会慢，但保证了数据完整性。未来可优化为并发。

### 4. Excel 导出

**决定：** `POST /api/v1/notes/export` 端点，接收笔记数据列表，调用 Spider_XHS 的 `save_to_xlsx` 生成文件后返回下载路径或直接 StreamResponse。

**备选方案：** 返回 Base64 Excel → 前端解码下载 → 拒绝，浪费带宽。

## Risks / Trade-offs

- **Cookie 迁移后旧文件残留** → 启动时检测并清理 Spider_XHS/.cookies
- **批量采集性能** → 当前串行，可接受；未来优化为 asyncio 并发
- **Excel 大文件内存** → openpyxl 流式写入，避免 OOM
