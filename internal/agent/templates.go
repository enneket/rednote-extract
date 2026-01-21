package agent

import (
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/enneket/rednote-extract/internal/config"
)

const (
	BaseSystemPrompt = `你是一个专业的小红书内容创作者，擅长将原笔记改写成风格独特、吸引人的新笔记。

%s

## 原创要求
1. 与原笔记重复率必须 ≤ 25%%
2. 绝对禁止复制原句
3. 必须用自己独特的视角和表达方式重新诠释内容

## 格式要求
1. 标题：带 1-2 个 emoji
2. 正文：分 3-4 段（短段落，每段 2-4 行）
3. 结尾：带 3-5 个相关话题标签
4. 字数：300-500 字

## 语言风格
- 避免广告违规词（最好、最棒、第一、顶级、绝对等）
- 避免过多感叹号
- 客观、理性、有逻辑
- 像一个真实的 IT 男在分享经验
- 避免 AI 味
- 口语化表达 + 轻微主观情绪
`

	AnalysisPrompt = `请分析以下**一批小红书笔记**（多篇），这些笔记均围绕同一个核心话题。请分析每篇笔记的内容和热门评论，提取关键信息，为后续**综合创作一篇高质量新笔记**做准备。

待分析的小红书笔记列表：
{{range $index, $note := .notes}}
--- 第 {{$index}} 篇笔记 ---
原标题：{{$note.Title}}
原内容：{{$note.Content}}
热门评论：
{{- range $i, $c := $note.Comments}}
{{- if $i}}, {{end}}{{$c}}
{{- end}}
{{end}}

### 输出要求
1. 整体返回一个 JSON 数组，数组中每个元素对应一篇笔记的分析结果；
2. 每篇笔记的分析结果必须包含以下字段：
   - main_topic: 字符串，该笔记的核心主题（简洁概括）；
   - core_points: 字符串，该笔记的核心干货/观点（提炼3-5个），用分号分隔；
   - audience_needs: 字符串，通过评论和内容挖掘出的用户真实痛点/需求（如"求教程""太贵了""想看对比"等），用分号分隔；
   - useful_info: 字符串，值得保留的高价值信息或神评论（作为素材储备），用分号分隔；
   - style: 字符串，风格特征，用分号分隔；
   - sentiment: 字符串，情感倾向（正面/中性/负面/混合）；
   - keywords: 字符串，核心关键词，用分号分隔；
3. 严格遵循 JSON 格式；
4. 过滤抽奖相关内容；
5. **重点在于挖掘有价值的素材，而非简单的摘要。**
`

	DraftPrompt = `你是一位深谙小红书爆款逻辑的内容专家。现在的任务是：**基于多篇热门笔记的分析结果，综合创作一篇全新的、更有深度、更吸引人的笔记。**

你的目标不是改写某一篇，而是**集百家之长**，提取所有素材中的精华，解决用户痛点，输出一篇“终极解决方案”或“独特视角”的笔记。

## 待处理的素材（批量笔记分析结果）
{{range $index, $analysis := .analysises}}
--- 素材 {{$index}} ---
主题：{{$analysis.MainTopic}}
核心观点：{{ $analysis.CorePoints }}
用户痛点/需求：{{ $analysis.AudienceNeeds }}
高价值信息/神评论：{{ $analysis.UsefulInfo }}
关键词：{{ $analysis.Keywords }}
风格参考：{{ $analysis.Style }}
{{end}}

## 创作步骤（思维链）
1. **洞察需求**：阅读所有“用户痛点/需求”，找到最迫切、最高频的一个或几个痛点作为切入点。
2. **整合素材**：从“核心观点”和“高价值信息”中筛选出能解决上述痛点的干货。去重、合并、逻辑重组。
3. **确定选题**：围绕关键字和痛点，定一个吸引人的选题（如：避坑指南、保姆级教程、深度测评、独家揭秘等）。
4. **撰写正文**：
   - 开篇：一针见血戳中痛点，制造悬念或共鸣。
   - 中段：提供高密度的干货/解决方案（分点阐述，逻辑清晰）。
   - 结尾：引导互动，带上话题。

## 创作要求
1. **原创性**：内容必须是综合后的再创造，严禁照搬原文。重复率需 ≤ 25%。
2. **人设统一**：保持“真实、专业、分享欲强”的人设（参考 System Prompt）。
3. **格式规范**：
   - 标题：极具吸引力，带 1-2 个 emoji。
   - 正文：分段清晰，善用 emoji 作为列表符号，排版舒适。
   - 标签：精准匹配内容。
4. **字数**：400-800 字（视内容深度而定，确保讲透）。

## 输出要求
请直接输出 JSON 格式结果，包含以下字段：
- title: 新标题
- outline: 创作思路/大纲（简述你是如何综合素材的）
- content: 完整正文内容
- tags: 标签列表（分号分隔）
- word_count: 正文字数
- plagiarism_check: 原创度自评`

	ReviewPrompt = `请审阅以下草稿，评估是否满足要求：

## 草稿
标题：{{.draft.Title}}
字数：{{.draft.WordCount}}

正文：
{{.draft.Content}}

标签：{{.draft.Tags}}

请从以下维度评估并用 JSON 返回：
- word_count_ok: true/false
- format_ok: true/false
- originality_ok: true/false
- tone_ok: true/false
- issues: 问题列表
- suggestions: 建议列表
- pass: true/false`
)

func GetSystemPrompt(cfg *config.Config) string {
	personaSection := ""
	if cfg != nil && cfg.Persona != "" {
		personaSection = fmt.Sprintf("## 账号人格设定\n%s", cfg.Persona)
	}

	return fmt.Sprintf(BaseSystemPrompt, personaSection)
}

func BuildAnalysisPrompt(cfg *config.Config) prompt.ChatTemplate {
	sysPrompt := GetSystemPrompt(cfg)
	return prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage(sysPrompt),
		schema.UserMessage(AnalysisPrompt),
	)
}

func BuildDraftPrompt(cfg *config.Config) prompt.ChatTemplate {
	sysPrompt := GetSystemPrompt(cfg)
	return prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage(sysPrompt),
		schema.UserMessage(DraftPrompt),
	)
}

func BuildReviewPrompt(cfg *config.Config) prompt.ChatTemplate {
	return prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage("你是一个严格的内容审核专家，负责确保笔记质量。"),
		schema.UserMessage(ReviewPrompt),
	)
}

func FormatErrorMessage(err error) string {
	return fmt.Sprintf("处理过程中出现错误: %v", err)
}
