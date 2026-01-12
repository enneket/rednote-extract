package utils

import (
	"testing"

	"github.com/enneket/rednote-extract/models"
)

func TestParseJSONWithCleanup(t *testing.T) {
	// 测试正常JSON
	t.Run("Valid JSON", func(t *testing.T) {
		jsonStr := `{
			"main_topic": "测试主题",
			"core_points": ["点1", "点2"],
			"audience_needs": ["需求1", "需求2"],
			"style": ["风格1", "风格2"],
			"sentiment": "正面",
			"keywords": ["关键词1", "关键词2"]
		}`

		var result models.AnalyzedInput
		err := ParseJSONWithCleanup(jsonStr, &result)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.MainTopic != "测试主题" {
			t.Errorf("Expected MainTopic to be '测试主题', got '%s'", result.MainTopic)
		}
	})

	// 测试带markdown格式的JSON
	t.Run("JSON with markdown formatting", func(t *testing.T) {
		jsonStr := "```" + `json
{
	"main_topic": "测试主题2",
	"core_points": ["点A", "点B"],
	"audience_needs": ["需求A", "需求B"],
	"style": ["风格A", "风格B"],
	"sentiment": "正面",
	"keywords": ["关键词A", "关键词B"]
}` + "```"

		var result models.AnalyzedInput
		err := ParseJSONWithCleanup(jsonStr, &result)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.MainTopic != "测试主题2" {
			t.Errorf("Expected MainTopic to be '测试主题2', got '%s'", result.MainTopic)
		}
	})
}
