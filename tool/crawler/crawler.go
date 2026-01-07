package crawler

import (
	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	innerCrawler "github.com/enneket/rednote-extract/tool/crawler/internal/crawler"
	"github.com/enneket/rednote-extract/tool/crawler/internal/model"
	"github.com/enneket/rednote-extract/tool/crawler/internal/store"
	"github.com/enneket/rednote-extract/tool/crawler/internal/tools"
)

// Crawler 对外暴露的爬虫接口
type Crawler struct {
	crawler *innerCrawler.RednoteCrawler
	config  *config.Config
}

// NewCrawler 创建爬虫实例
func NewCrawler(cookies string) (*Crawler, error) {
	// 加载配置
	cfg := config.DefaultConfig()
	cfg.Cookies = cookies

	// 创建日志记录器
	logger := tools.NewLogger()

	// 创建存储
	store, err := store.NewJSONStore(cfg)
	if err != nil {
		return nil, err
	}

	// 创建内部爬虫实例
	innerCrawler := innerCrawler.NewRednoteCrawler(cfg, store, logger)

	// 初始化爬虫
	if err := innerCrawler.Init(); err != nil {
		return nil, err
	}

	return &Crawler{
		crawler: innerCrawler,
		config:  cfg,
	}, nil
}

// Search 搜索帖子
// keyword: 搜索关键词
// return: 帖子列表和错误信息
func (c *Crawler) Search(keyword string) ([]*model.Note, error) {
	return c.crawler.Search(keyword)
}

// GetNoteDetail 获取帖子详情
// noteURL: 帖子URL
// return: 帖子详情和错误信息
func (c *Crawler) GetNoteDetail(noteURL string) (*model.Note, error) {
	return c.crawler.GetSpecifiedNotes(noteURL)
}

// GetNoteComments 获取帖子评论
// noteID: 帖子ID
// xsecToken: 安全令牌
// return: 评论列表和错误信息
func (c *Crawler) GetNoteComments(noteID, xsecToken string) ([]*model.Comment, error) {
	return c.crawler.GetNoteComments(noteID, xsecToken)
}

// Close 关闭爬虫
func (c *Crawler) Close() error {
	return c.crawler.Close()
}
