## ADDED Requirements

### Requirement: Excel export endpoint
The Gateway SHALL expose POST /api/v1/notes/export that accepts note data and returns an Excel file download.

#### Scenario: Export notes to Excel
- **WHEN** POST /api/v1/notes/export with array of note objects
- **THEN** Gateway generates an .xlsx file and returns it as a downloadable attachment with Content-Disposition header

#### Scenario: Export with custom filename
- **WHEN** POST /api/v1/notes/export with `filename` parameter
- **THEN** the downloaded file uses that name (e.g., `榴莲笔记.xlsx`)

#### Scenario: Empty notes array
- **WHEN** POST /api/v1/notes/export with empty array
- **THEN** Gateway returns 400 error with `{"error": "no notes to export"}`

### Requirement: Excel format requirements
The exported Excel SHALL contain columns matching the unified note structure.

#### Scenario: Excel columns
- **WHEN** notes are exported
- **THEN** Excel contains columns: ID, 标题, 作者, 作者ID, 类型, 点赞数, 收藏数, 评论数, 分享数, URL, 标签, 发布时间, IP属地
