package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LLMProvider   string  `json:"llm_provider"`
	APIBaseURL    string  `json:"api_base_url"`
	APIKey        string  `json:"api_key"`
	ModelName     string  `json:"model_name"`
	Temperature   float32 `json:"temperature"`
	MaxIterations int     `json:"max_iterations"`
}

var globalConfig *Config

func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	log.Printf("DEBUG: API_KEY from env: %s", apiKey)
	if apiKey == "" {
		apiKey = "sk-placeholder-key"
		log.Printf("DEBUG: Using placeholder API key")
	}

	baseURL := os.Getenv("API_BASE_URL")
	log.Printf("DEBUG: API_BASE_URL from env: %s", baseURL)
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	llmProvider := getEnv("LLM_PROVIDER", "openai")
	log.Printf("DEBUG: LLM_PROVIDER from env: %s", llmProvider)

	modelName := getEnv("MODEL_NAME", "gpt-4o")
	log.Printf("DEBUG: MODEL_NAME from env: %s", modelName)

	maxIterations, err := strconv.Atoi(getEnv("MAX_ITERATIONS", "5"))
	if err != nil {
		log.Printf("Warning: Failed to parse MAX_ITERATIONS: %v", err)
		maxIterations = 5
	}
	globalConfig = &Config{
		LLMProvider:   llmProvider,
		APIBaseURL:    baseURL,
		APIKey:        apiKey,
		ModelName:     modelName,
		Temperature:   0.7,
		MaxIterations: maxIterations,
	}

	log.Printf("Config loaded!")

	return globalConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var ctx = context.Background()
