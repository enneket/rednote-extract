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

	AnalysisPrompt = `请分析以下**一批小红书笔记**（多篇），为每一篇笔记独立提取关键信息（用于后续改写），最终用 JSON 格式返回所有笔记的分析结果。

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
   - main_topic: 字符串，该笔记的核心主题（简洁概括，不超过20字）；
   - core_points: 字符串，该笔记的核心要点列表（每个要点为字符串，提炼3-5个核心信息），用分号分隔；
   - audience_needs: 字符串，读者可能的需求列表（基于笔记内容+评论推导，如"求教程""想避坑""找平替"等），用分号分隔；
   - style: 字符串，该笔记的风格特征（如"口语化""干货型""情绪化""图文感""种草向"等），用分号分隔；
   - sentiment: 字符串，情感倾向（可选值：正面/中性/负面/混合）；
   - keywords: 字符串，该笔记的核心关键词列表（提取5-8个，包含产品/场景/情绪等维度），用分号分隔；
3. 严格遵循 JSON 格式，确保无语法错误，字段值类型符合要求；
4. 分析需贴合每篇笔记的实际内容，不同笔记的分析结果需独立且精准，不混淆信息;
5. 需要过滤抽奖的相关笔记和评论，避免在笔记中包含抽奖信息，这十分重要!
`

	DraftPrompt = `基于以下批量小红书笔记的分析结果，为每一篇笔记独立改写一篇全新的小红书笔记。
## 批量笔记分析结果列表
{{range $index, $analysis := .analysises}}
--- 第 {{$index}} 篇笔记分析结果 ---
主题：{{$analysis.MainTopic}}
核心要点：{{ $analysis.CorePoints }}
读者需求：{{ $analysis.AudienceNeeds }}
关键词：{{ $analysis.Keywords }}
风格（参考）：{{ $analysis.Style }}
情感倾向（参考）：{{$analysis.Sentiment}}
{{end}}

## 改写核心要求
1.  视角与风格：严格以人设视角创作，语言风格需贴合该群体特征，避免过度柔美或浮夸表达；
2.  内容要求：每篇改写笔记需完整覆盖对应分析结果的核心要点，同时精准匹配读者需求（如读者需“避坑”则重点补充实用避坑技巧，需“教程”则强化步骤清晰度），不得遗漏关键信息；
3.  格式规范：每篇笔记标题必须带1-2个贴合主题的emoji，正文分3-4段（每段聚焦1个核心要点，逻辑连贯），结尾需搭配适配的话题标签;
4.  权重判定不以文本字数多少为依据，优先依据字符出现的频率（字频）进行内容改写。


## 输出要求
输出一篇笔记，请用 JSON 格式返回，包含以下字段：
- title: 字符串，带 emoji 的新标题
- outline: 字符串，文章结构说明
- content: 字符串，完整正文内容
- tags: 字符串，标签列表（每个标签为字符串，用分号分隔）
- word_count: 字数（整数）
- plagiarism_check: 原创度自评（高/中/低）`

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
