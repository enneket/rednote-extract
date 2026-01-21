package llm

import (
	"context"
	"fmt"

	einoOpenAI "github.com/cloudwego/eino-ext/components/model/openai"
	einoQwen "github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/enneket/rednote-extract/internal/config"
)

type ChatModel interface {
	Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error)
}

func NewChatModel(cfg *config.Config) (ChatModel, error) {
	var model model.ToolCallingChatModel
	var err error

	temperature := cfg.Temperature

	switch cfg.LLMProvider {
	case "openai":
		model, err = einoOpenAI.NewChatModel(context.TODO(), &einoOpenAI.ChatModelConfig{
			BaseURL:     cfg.LLMAPIBaseURL,
			APIKey:      cfg.LLMAPIKey,
			Model:       cfg.ModelName,
			Temperature: &temperature,
		})
	case "qwen":
		model, err = einoQwen.NewChatModel(context.TODO(), &einoQwen.ChatModelConfig{
			BaseURL:     cfg.LLMAPIBaseURL,
			APIKey:      cfg.LLMAPIKey,
			Model:       cfg.ModelName,
			Temperature: &temperature,
		})
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLMProvider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create chat model: %w", err)
	}

	return model, nil
}
