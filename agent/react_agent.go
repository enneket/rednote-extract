package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/enneket/rednote-extract/chat_model"
	"github.com/enneket/rednote-extract/config"
	"github.com/enneket/rednote-extract/models"
	"github.com/enneket/rednote-extract/prompt"
	"github.com/enneket/rednote-extract/utils"
)

type ReactAgent struct {
	cfg   *config.Config
	model chat_model.ChatModel
}

func NewReactAgent(cfg *config.Config) *ReactAgent {
	chatModel, err := chat_model.NewChatModel(cfg)
	if err != nil {
		log.Fatalf("Failed to create chat model: %v", err)
	}
	return &ReactAgent{
		cfg:   cfg,
		model: chatModel,
	}
}

func (a *ReactAgent) GenerateNote(ctx context.Context, input []*models.NoteInput) (*models.GeneratedNote, error) {
	state := &models.AgentState{
		Input:          input,
		ThoughtProcess: make([]string, 0),
		Iteration:      0,
	}

	log.Printf("[DEBUG] GenerateNote: 开始处理笔记")

	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 开始处理笔记", time.Now().Format("15:04:05")))

	state, err := a.analyzeNode(ctx, state)
	if err != nil {
		log.Printf("[DEBUG] GenerateNote: analyzeNode 失败: %v", err)
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	state, err = a.draftNode(ctx, state)
	if err != nil {
		log.Printf("[DEBUG] GenerateNote: draftNode 失败: %v", err)
		return nil, fmt.Errorf("draft failed: %w", err)
	}

	iteration := 0
	for state.DraftNote != nil && !state.DraftNote.IsValid() && state.Iteration < a.cfg.MaxIterations {
		state.Iteration++
		iteration++
		log.Printf("[DEBUG] GenerateNote: 开始第 %d 次 review-revision 循环, iteration=%d", iteration, state.Iteration)

		state, err = a.reviewNode(ctx, state)
		if err != nil {
			log.Printf("[DEBUG] GenerateNote: reviewNode 失败: %v", err)
			return nil, fmt.Errorf("review failed: %w", err)
		}

		log.Printf("[DEBUG] GenerateNote: reviewNode 完成, pass=%v, issues_count=%d",
			state.DraftNote.ReviewResult.Pass, len(state.DraftNote.ReviewResult.Issues))

		state, err = a.reviseNode(ctx, state)
		if err != nil {
			log.Printf("[DEBUG] GenerateNote: reviseNode 失败: %v", err)
			return nil, fmt.Errorf("revision failed: %w", err)
		}

		log.Printf("[DEBUG] GenerateNote: reviseNode 完成, draft_word_count=%d", state.DraftNote.WordCount)
	}

	log.Printf("[DEBUG] GenerateNote: 开始 final review")
	state, err = a.reviewNode(ctx, state)
	if err != nil {
		log.Printf("[DEBUG] GenerateNote: final reviewNode 失败: %v", err)
		return nil, fmt.Errorf("final review failed: %w", err)
	}

	state, err = a.finalizeNode(ctx, state)
	if err != nil {
		log.Printf("[DEBUG] GenerateNote: finalizeNode 失败: %v", err)
		return nil, fmt.Errorf("finalize failed: %w", err)
	}

	log.Printf("[DEBUG] GenerateNote: 完成, final_title=%s, final_word_count=%d",
		state.FinalNote.Title, len(state.FinalNote.Content))

	return state.FinalNote, nil
}

func (a *ReactAgent) analyzeNode(ctx context.Context, state *models.AgentState) (*models.AgentState, error) {
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 开始分析原笔记...", time.Now().Format("15:04:05")))

	log.Printf("[DEBUG] analyzeNode: 开始分析原笔记...")

	promptTpl := prompt.BuildAnalysisPrompt()
	messages, err := promptTpl.Format(ctx, map[string]any{
		"notes": state.Input,
	})
	if err != nil {
		log.Printf("[DEBUG] analyzeNode: format prompt 失败: %v", err)
		return nil, fmt.Errorf("failed to format analysis prompt: %w", err)
	}

	for _, msg := range messages {
		log.Printf("[DEBUG] analyzeNode: message: %s", msg.Content)
	}

	log.Printf("[DEBUG] analyzeNode: 调用 model 生成分析结果")
	response, err := a.model.Generate(ctx, messages)
	if err != nil {
		log.Printf("[DEBUG] analyzeNode: model.Generate 失败: %v", err)
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	log.Printf("[DEBUG] analyzeNode: model 响应: %s", response.Content)

	var analyzed []*models.AnalyzedInput
	if err := utils.ParseJSONWithCleanup(response.Content, &analyzed); err != nil {
		log.Printf("[DEBUG] analyzeNode: JSON 解析失败: %v", err)
		return nil, fmt.Errorf("failed to parse analysis result: %w", err)
	}

	state.AnalyzedInput = analyzed
	log.Printf("[DEBUG] analyzeNode: 完成")
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 分析完成", time.Now().Format("15:04:05")))

	return state, nil
}

func (a *ReactAgent) draftNode(ctx context.Context, state *models.AgentState) (*models.AgentState, error) {
	state.Iteration++
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 开始第 %d 次草稿生成...", time.Now().Format("15:04:05"), state.Iteration))

	log.Printf("[DEBUG] draftNode: 开始第 %d 次草稿生成", state.Iteration)

	if state.AnalyzedInput == nil {
		log.Printf("[DEBUG] draftNode: 错误 - AnalyzedInput 为 nil")
		return nil, fmt.Errorf("missing analysis result")
	}

	promptTpl := prompt.BuildDraftPrompt()
	messages, err := promptTpl.Format(ctx, map[string]any{
		"analysises": state.AnalyzedInput,
	})
	if err != nil {
		log.Printf("[DEBUG] draftNode: format prompt 失败: %v", err)
		return nil, fmt.Errorf("failed to format draft prompt: %w", err)
	}
	for _, msg := range messages {
		log.Printf("[DEBUG] draftNode: message: %s", msg.Content)
	}

	log.Printf("[DEBUG] draftNode: 调用 model 生成草稿")
	response, err := a.model.Generate(ctx, messages)
	if err != nil {
		log.Printf("[DEBUG] draftNode: model.Generate 失败: %v", err)
		return nil, fmt.Errorf("draft generation failed: %w", err)
	}

	log.Printf("[DEBUG] draftNode: model 响应: %s", response.Content)

	var draft models.DraftNote
	if err := utils.ParseJSONWithCleanup(response.Content, &draft); err != nil {
		log.Printf("[DEBUG] draftNode: JSON 解析失败: %v", err)
		return nil, fmt.Errorf("failed to parse draft result: %w", err)
	}

	draft.WordCount = len(draft.Content)
	state.DraftNote = &draft
	log.Printf("[DEBUG] draftNode: 完成")
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 草稿生成完成，字数: %d", time.Now().Format("15:04:05"), draft.WordCount))

	return state, nil
}

func (a *ReactAgent) reviewNode(ctx context.Context, state *models.AgentState) (*models.AgentState, error) {
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 开始审核草稿...", time.Now().Format("15:04:05")))

	log.Printf("[DEBUG] reviewNode: 开始审核草稿")

	if state.DraftNote == nil {
		log.Printf("[DEBUG] reviewNode: 错误 - DraftNote 为 nil")
		return nil, fmt.Errorf("missing draft note")
	}

	log.Printf("[DEBUG] reviewNode: 待审核草稿, title=%s, word_count=%d, tags=%v",
		state.DraftNote.Title, state.DraftNote.WordCount, state.DraftNote.Tags)

	promptTpl := prompt.BuildReviewPrompt()
	messages, err := promptTpl.Format(ctx, map[string]any{
		"draft": state.DraftNote,
	})
	if err != nil {
		log.Printf("[DEBUG] reviewNode: format prompt 失败: %v", err)
		return nil, fmt.Errorf("failed to format review prompt: %w", err)
	}

	log.Printf("[DEBUG] reviewNode: 调用 model 进行审核")
	response, err := a.model.Generate(ctx, messages)
	if err != nil {
		log.Printf("[DEBUG] reviewNode: model.Generate 失败: %v", err)
		return nil, fmt.Errorf("review failed: %w", err)
	}

	log.Printf("[DEBUG] reviewNode: model 响应: %s", response.Content)

	var review *models.ReviewResult
	if err := utils.ParseJSONWithCleanup(response.Content, &review); err != nil {
		log.Printf("[DEBUG] reviewNode: JSON 解析失败，使用自动评估: %v", err)
		state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 审核解析失败，使用自动评估", time.Now().Format("15:04:05")))
		review = autoReview(state.DraftNote)
	} else {
		log.Printf("[DEBUG] reviewNode: 审核解析成功")
	}

	state.DraftNote.ReviewResult = review
	log.Printf("[DEBUG] reviewNode: 完成, pass=%v, word_count_ok=%v, format_ok=%v, originality_ok=%v, tone_ok=%v, issues_count=%d",
		review.Pass, review.WordCountOK, review.FormatOK, review.OriginalityOK, review.ToneOK, len(review.Issues))
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 审核完成，通过: %v", time.Now().Format("15:04:05"), review.Pass))

	return state, nil
}

func (a *ReactAgent) reviseNode(ctx context.Context, state *models.AgentState) (*models.AgentState, error) {
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 开始修订...", time.Now().Format("15:04:05")))

	log.Printf("[DEBUG] reviseNode: 开始修订")

	if state.DraftNote == nil || state.DraftNote.ReviewResult.Issues == nil {
		log.Printf("[DEBUG] reviseNode: 错误 - DraftNote 或 ReviewResult.Issues 为 nil")
		return nil, fmt.Errorf("nothing to revise")
	}

	log.Printf("[DEBUG] reviseNode: 待修订草稿, title=%s, word_count=%d, issues_count=%d",
		state.DraftNote.Title, state.DraftNote.WordCount, len(state.DraftNote.ReviewResult.Issues))
	log.Printf("[DEBUG] reviseNode: 问题: %v", state.DraftNote.ReviewResult.Issues)
	log.Printf("[DEBUG] reviseNode: 建议: %v", state.DraftNote.ReviewResult.Suggestions)

	revisionPrompt := fmt.Sprintf(`请根据以下审核意见改进笔记：

## 原草稿
标题：%s
字数：%d

正文：
%s

标签：%s

## 审核意见
问题：%v
建议：%v

请生成改进后的笔记，确保解决所有问题。`, state.DraftNote.Title, state.DraftNote.WordCount, state.DraftNote.Content, state.DraftNote.Tags, state.DraftNote.ReviewResult.Issues, state.DraftNote.ReviewResult.Suggestions)

	messages := []*schema.Message{
		schema.SystemMessage(prompt.SystemPrompt),
		schema.UserMessage(revisionPrompt),
	}

	log.Printf("[DEBUG] reviseNode: 调用 model 进行修订")
	response, err := a.model.Generate(ctx, messages)
	if err != nil {
		log.Printf("[DEBUG] reviseNode: model.Generate 失败: %v", err)
		return nil, fmt.Errorf("revision failed: %w", err)
	}

	log.Printf("[DEBUG] reviseNode: model 响应长度=%d", len(response.Content))

	var revised models.DraftNote
	if err := utils.ParseJSONWithCleanup(response.Content, &revised); err != nil {
		log.Printf("[DEBUG] reviseNode: JSON 解析失败，保留原草稿: %v", err)
		state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 修订解析失败，保留原草稿", time.Now().Format("15:04:05")))
		return state, nil
	}

	revised.WordCount = len(revised.Content)
	state.DraftNote = &revised

	log.Printf("[DEBUG] reviseNode: 完成, title=%s, word_count=%d", revised.Title, revised.WordCount)
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 修订完成", time.Now().Format("15:04:05")))

	return state, nil
}

func (a *ReactAgent) finalizeNode(ctx context.Context, state *models.AgentState) (*models.AgentState, error) {
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 完成最终处理", time.Now().Format("15:04:05")))

	log.Printf("[DEBUG] finalizeNode: 开始最终处理")

	if state.DraftNote == nil {
		log.Printf("[DEBUG] finalizeNode: 错误 - DraftNote 为 nil")
		return nil, fmt.Errorf("no draft to finalize")
	}

	log.Printf("[DEBUG] finalizeNode: 待处理草稿, title=%s, word_count=%d, tags=%v",
		state.DraftNote.Title, state.DraftNote.WordCount, state.DraftNote.Tags)

	state.FinalNote = &models.GeneratedNote{
		Title:   state.DraftNote.Title,
		Content: state.DraftNote.Content,
		Tags:    strings.Join(strings.Split(state.DraftNote.Tags, ";"), " "),
	}

	log.Printf("[DEBUG] finalizeNode: 完成, final_title=%s, final_content_length=%d, final_tags=%s",
		state.FinalNote.Title, len(state.FinalNote.Content), state.FinalNote.Tags)
	state.ThoughtProcess = append(state.ThoughtProcess, fmt.Sprintf("[%s] 笔记生成完成", time.Now().Format("15:04:05")))

	return state, nil
}

func autoReview(draft *models.DraftNote) *models.ReviewResult {
	result := &models.ReviewResult{
		WordCountOK:   draft.WordCount >= models.MinWordCount && draft.WordCount <= models.MaxWordCount,
		FormatOK:      true,
		OriginalityOK: true,
		ToneOK:        true,
		Issues:        make([]string, 0),
		Suggestions:   make([]string, 0),
		Pass:          true,
	}

	if !result.WordCountOK {
		result.Issues = append(result.Issues, fmt.Sprintf("字数不在 %d-%d 范围内", models.MinWordCount, models.MaxWordCount))
		result.Suggestions = append(result.Suggestions, "调整内容长度")
		result.Pass = false
	}

	if len(draft.Tags) < 3 {
		result.Issues = append(result.Issues, "话题标签不足 3 个")
		result.Suggestions = append(result.Suggestions, "添加更多相关话题标签")
		result.Pass = false
	}

	return result
}
