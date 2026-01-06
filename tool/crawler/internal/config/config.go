package config

// Config 应用配置结构体，定义了爬虫的所有配置项
type Config struct {
	// Keyword 搜索关键词，用于爬取指定关键词的内容
	Keyword string `mapstructure:"KEYWORD"`
	// Cookies 登录Cookie字符串，当LoginType为cookie时使用
	Cookies string `mapstructure:"COOKIES"`
	// Headless 浏览器是否以无头模式运行
	Headless bool `mapstructure:"HEADLESS"`
	// SaveLoginState 是否保存登录状态到本地
	SaveLoginState bool `mapstructure:"SAVE_LOGIN_STATE"`
	// CustomBrowserPath 自定义浏览器路径
	CustomBrowserPath string `mapstructure:"CUSTOM_BROWSER_PATH"`
	// BrowserLaunchTimeout 浏览器启动超时时间（秒）
	BrowserLaunchTimeout int `mapstructure:"BROWSER_LAUNCH_TIMEOUT"`
	// AutoCloseBrowser 程序退出时是否自动关闭浏览器
	AutoCloseBrowser bool `mapstructure:"AUTO_CLOSE_BROWSER"`
	// SaveDataOption 数据保存方式，支持 "json"（JSON文件）
	SaveDataOption string `mapstructure:"SAVE_DATA_OPTION"`
	// UserDataDir 用户数据目录路径模板，%s会被替换为平台名称
	UserDataDir string `mapstructure:"USER_DATA_DIR"`
	// StartPage 起始页码，从第几页开始爬取
	StartPage int `mapstructure:"START_PAGE"`
	// CrawlerMaxNotesCount 单次爬取的最大笔记数量
	CrawlerMaxNotesCount int `mapstructure:"CRAWLER_MAX_NOTES_COUNT"`
	// MaxConcurrencyNum 最大并发数
	MaxConcurrencyNum int `mapstructure:"MAX_CONCURRENCY_NUM"`
	// EnableGetComments 是否爬取笔记的评论
	EnableGetComments bool `mapstructure:"ENABLE_GET_COMMENTS"`
	// CrawlerMaxCommentsCountSingleNotes 单条笔记的最大评论爬取数量
	CrawlerMaxCommentsCountSingleNotes int `mapstructure:"CRAWLER_MAX_COMMENTS_COUNT_SINGLENOTES"`
	// EnableGetSubComments 是否爬取评论的子评论（回复）
	EnableGetSubComments bool `mapstructure:"ENABLE_GET_SUB_COMMENTS"`
	// CrawlerMaxSleepSec 爬虫最大休眠时间（秒），用于控制请求频率
	CrawlerMaxSleepSec int `mapstructure:"CRAWLER_MAX_SLEEP_SEC"`
	// SortType 排序方式，如 "time"（时间排序）、"popular"（热度排序）
	SortType string `mapstructure:"SORT_TYPE"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Cookies:                            "",
		Headless:                           false,
		SaveLoginState:                     true,
		CustomBrowserPath:                  "",
		BrowserLaunchTimeout:               60,
		AutoCloseBrowser:                   true,
		SaveDataOption:                     "json",
		UserDataDir:                        "%s_user_data_dir",
		StartPage:                          1,
		CrawlerMaxNotesCount:               50,
		MaxConcurrencyNum:                  1,
		EnableGetComments:                  true,
		CrawlerMaxCommentsCountSingleNotes: 50,
		EnableGetSubComments:               true,
		CrawlerMaxSleepSec:                 2,
		SortType:                           "popularity_descending",
	}
}
