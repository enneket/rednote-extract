package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// LLM Config
	LLMProvider   string  `mapstructure:"LLM_PROVIDER"`
	LLMAPIBaseURL string  `mapstructure:"LLM_API_BASE_URL"`
	LLMAPIKey     string  `mapstructure:"LLM_API_KEY"`
	ModelName     string  `mapstructure:"MODEL_NAME"`
	Temperature   float32 `mapstructure:"TEMPERATURE"`
	MaxIterations int     `mapstructure:"MAX_ITERATIONS"`

	// Crawler Config
	Keywords           string `mapstructure:"KEYWORDS"`
	Cookies            string `mapstructure:"COOKIES"`
	Headless           bool   `mapstructure:"HEADLESS"`
	SaveDataOption     string `mapstructure:"SAVE_DATA_OPTION"`
	UserDataDir        string `mapstructure:"USER_DATA_DIR"`
	EnableGetComments  bool   `mapstructure:"ENABLE_GET_COMMENTS"`
	CrawlerMaxSleepSec int    `mapstructure:"CRAWLER_MAX_SLEEP_SEC"`
	MaxNotes           int    `mapstructure:"MAX_NOTES"`

	// Prompts
	Persona string `mapstructure:"PERSONA"`
}

var AppConfig *Config

func Load() (*Config, error) {
	if AppConfig != nil {
		return AppConfig, nil
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Set Defaults
	// LLM Defaults
	viper.SetDefault("LLM_PROVIDER", "openai")
	viper.SetDefault("LLM_API_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("MODEL_NAME", "gpt-4o")
	viper.SetDefault("TEMPERATURE", 0.7)
	viper.SetDefault("MAX_ITERATIONS", 5)

	// Crawler Defaults
	viper.SetDefault("HEADLESS", false)
	viper.SetDefault("SAVE_DATA_OPTION", "json")
	viper.SetDefault("ENABLE_GET_COMMENTS", true)
	viper.SetDefault("CRAWLER_MAX_SLEEP_SEC", 2)
	viper.SetDefault("MAX_NOTES", 10)

	// Also allow env vars with prefix or raw
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	AppConfig = config
	return AppConfig, nil
}

func GetKeywords() []string {
	if AppConfig == nil || AppConfig.Keywords == "" {
		return []string{}
	}
	return strings.Split(AppConfig.Keywords, ",")
}
