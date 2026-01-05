package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader 配置加载器
type Loader struct {
	viper *viper.Viper
}

// NewLoader 创建配置加载器
func NewLoader() *Loader {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")

	// 支持从环境变量加载配置
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Loader{
		viper: v,
	}
}

// Load 加载配置
func (l *Loader) Load() (*Config, error) {
	// 读取配置文件
	if err := l.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在，使用默认配置
		fmt.Println("config file not found, using default config")
	}

	// 设置默认值
	l.setDefaults()

	// 解析配置
	config := DefaultConfig()
	if err := l.viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 特殊处理关键词
	if keywordsStr := l.viper.GetString("KEYWORDS"); keywordsStr != "" {
		config.Keywords = ParseKeywords(keywordsStr)
	}

	return config, nil
}

// setDefaults 设置默认值
func (l *Loader) setDefaults() {
	defaultConfig := DefaultConfig()

	l.viper.SetDefault("KEYWORDS", strings.Join(defaultConfig.Keywords, ","))
	l.viper.SetDefault("LOGIN_TYPE", defaultConfig.LoginType)
	l.viper.SetDefault("COOKIES", defaultConfig.Cookies)
	l.viper.SetDefault("CRAWLER_TYPE", defaultConfig.CrawlerType)
	l.viper.SetDefault("HEADLESS", defaultConfig.Headless)
	l.viper.SetDefault("SAVE_LOGIN_STATE", defaultConfig.SaveLoginState)
	l.viper.SetDefault("CUSTOM_BROWSER_PATH", defaultConfig.CustomBrowserPath)
	l.viper.SetDefault("BROWSER_LAUNCH_TIMEOUT", defaultConfig.BrowserLaunchTimeout)
	l.viper.SetDefault("AUTO_CLOSE_BROWSER", defaultConfig.AutoCloseBrowser)
	l.viper.SetDefault("SAVE_DATA_OPTION", defaultConfig.SaveDataOption)
	l.viper.SetDefault("USER_DATA_DIR", defaultConfig.UserDataDir)
	l.viper.SetDefault("START_PAGE", defaultConfig.StartPage)
	l.viper.SetDefault("CRAWLER_MAX_NOTES_COUNT", defaultConfig.CrawlerMaxNotesCount)
	l.viper.SetDefault("MAX_CONCURRENCY_NUM", defaultConfig.MaxConcurrencyNum)
	l.viper.SetDefault("ENABLE_GET_MEIDAS", defaultConfig.EnableGetMedias)
	l.viper.SetDefault("ENABLE_GET_COMMENTS", defaultConfig.EnableGetComments)
	l.viper.SetDefault("CRAWLER_MAX_COMMENTS_COUNT_SINGLENOTES", defaultConfig.CrawlerMaxCommentsCountSingleNotes)
	l.viper.SetDefault("ENABLE_GET_SUB_COMMENTS", defaultConfig.EnableGetSubComments)
	l.viper.SetDefault("CRAWLER_MAX_SLEEP_SEC", defaultConfig.CrawlerMaxSleepSec)
	l.viper.SetDefault("SORT_TYPE", defaultConfig.SortType)
}
