package xhs

import (
	"fmt"
	"sync"
	"time"

	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	"github.com/enneket/rednote-extract/tool/crawler/internal/store"
	"github.com/enneket/rednote-extract/tool/crawler/internal/tools"
	"github.com/go-rod/rod"
)

// XiaoHongShuCrawler 小红书爬虫
type XiaoHongShuCrawler struct {
	config     *config.Config
	client     *XiaoHongShuClient
	browser    *rod.Browser
	page       *rod.Page
	store      store.Store
	browserMgr *tools.BrowserManager
	logger     tools.Logger
	wg         sync.WaitGroup
}

// NewXiaoHongShuCrawler 创建小红书爬虫
func NewXiaoHongShuCrawler(cfg *config.Config, store store.Store, logger tools.Logger) *XiaoHongShuCrawler {
	return &XiaoHongShuCrawler{
		config: cfg,
		store:  store,
		logger: logger,
	}
}

// Start 启动爬虫
func (c *XiaoHongShuCrawler) Start() error {
	c.logger.Info("[XiaoHongShuCrawler] Starting crawler")

	// 初始化浏览器管理器
	c.browserMgr = tools.NewBrowserManager(c.config, c.logger)

	// 启动浏览器
	browser, err := c.browserMgr.Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}
	c.browser = browser

	// 创建新页面
	page, err := c.browserMgr.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	c.page = page

	// 导航到小红书首页
	if err := page.Navigate("https://www.xiaohongshu.com").Error(); err != nil {
		return fmt.Errorf("failed to navigate to xiaohongshu: %w", err)
	}

	// 获取cookies
	cookies, err := c.browserMgr.GetCookies(".xiaohongshu.com")
	if err != nil {
		c.logger.Error("[XiaoHongShuCrawler] Failed to get cookies: %v", err)
	}
	cookieStr := c.browserMgr.ConvertCookiesToString(cookies)

	// 创建客户端
	c.client = NewXiaoHongShuClient(c.config, page, cookieStr, c.logger)

	// 测试客户端连接
	if !c.client.Pong() {
		// 登录
		if err := c.login(); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}
	}

	// 根据爬虫类型执行不同的任务
	switch c.config.CrawlerType {
	case "search":
		return c.search()
	case "detail":
		return c.getSpecifiedNotes()
	case "creator":
		return c.getCreatorsAndNotes()
	default:
		return fmt.Errorf("unsupported crawler type: %s", c.config.CrawlerType)
	}
}

// login 登录小红书
func (c *XiaoHongShuCrawler) login() error {
	c.logger.Info("[XiaoHongShuCrawler] Login required")

	// 这里实现登录逻辑
	// 支持扫码登录、手机登录等

	loginObj := NewXiaoHongShuLogin(c.config.LoginType, c.page, c.logger)
	if err := loginObj.Begin(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// 更新cookies
	cookies, err := c.browserMgr.GetCookies(".xiaohongshu.com")
	if err != nil {
		return fmt.Errorf("failed to get cookies after login: %w", err)
	}
	cookieStr := c.browserMgr.ConvertCookiesToString(cookies)
	c.client.UpdateCookies(cookieStr)

	return nil
}

// search 搜索帖子
func (c *XiaoHongShuCrawler) search() error {
	c.logger.Info("[XiaoHongShuCrawler] Begin search Xiaohongshu keywords")

	for _, keyword := range c.config.Keywords {
		c.logger.Info("[XiaoHongShuCrawler] Current search keyword: %s", keyword)

		page := 1
		searchID := tools.GetRandomString(16)

		for (page-c.config.StartPage+1)*20 <= c.config.CrawlerMaxNotesCount {
			if page < c.config.StartPage {
				c.logger.Info("[XiaoHongShuCrawler] Skip page %d", page)
				page++
				continue
			}

			c.logger.Info("[XiaoHongShuCrawler] Searching keyword: %s, page: %d", keyword, page)

			// 获取搜索结果
			notesRes, err := c.client.GetNoteByKeyword(keyword, searchID, page, c.config.SortType)
			if err != nil {
				c.logger.Error("[XiaoHongShuCrawler] Failed to get search results: %v", err)
				break
			}

			// 处理搜索结果
			// 这里需要根据实际API返回结构处理

			page++

			// 休眠
			time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
		}
	}

	return nil
}

// getSpecifiedNotes 获取指定帖子详情
func (c *XiaoHongShuCrawler) getSpecifiedNotes() error {
	c.logger.Info("[XiaoHongShuCrawler] Begin get specified notes")

	// 处理指定帖子URL列表
	for _, noteURL := range c.config.XHSSpecifiedNoteURLList {
		c.logger.Info("[XiaoHongShuCrawler] Processing note URL: %s", noteURL)

		// 解析帖子URL
		noteInfo := parseNoteInfoFromNoteURL(noteURL)
		if noteInfo == nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to parse note URL: %s", noteURL)
			continue
		}

		// 获取帖子详情
		note, err := c.getNoteDetail(noteInfo.NoteID, noteInfo.XsecSource, noteInfo.XsecToken)
		if err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to get note detail: %v", err)
			continue
		}

		// 保存帖子
		if err := c.store.SaveNote(note); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to save note: %v", err)
		}

		// 获取评论
		if c.config.EnableGetComments {
			if err := c.getNoteComments(note.NoteID, note.XsecToken); err != nil {
				c.logger.Error("[XiaoHongShuCrawler] Failed to get comments: %v", err)
			}
		}

		// 休眠
		time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
	}

	return nil
}

// getCreatorsAndNotes 获取创作者信息和帖子
func (c *XiaoHongShuCrawler) getCreatorsAndNotes() error {
	c.logger.Info("[XiaoHongShuCrawler] Begin get creators and notes")

	for _, creatorURL := range c.config.XHSCreatorIDList {
		c.logger.Info("[XiaoHongShuCrawler] Processing creator URL: %s", creatorURL)

		// 解析创作者URL
		creatorInfo := parseCreatorInfoFromURL(creatorURL)
		if creatorInfo == nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to parse creator URL: %s", creatorURL)
			continue
		}

		// 获取创作者信息
		creator, err := c.client.GetCreatorInfo(creatorInfo.UserID, creatorInfo.XsecToken, creatorInfo.XsecSource)
		if err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to get creator info: %v", err)
			continue
		}

		// 保存创作者信息
		if err := c.store.SaveCreator(creator); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to save creator: %v", err)
		}

		// 获取创作者的所有帖子
		allNotes, err := c.client.GetAllNotesByCreator(
			creatorInfo.UserID,
			c.config.CrawlerMaxSleepSec,
			c.processCreatorNotes,
			creatorInfo.XsecToken,
			creatorInfo.XsecSource,
		)
		if err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to get creator notes: %v", err)
			continue
		}

		// 批量获取评论
		if c.config.EnableGetComments {
			c.batchGetNoteComments(allNotes)
		}

		// 休眠
		time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
	}

	return nil
}

// getNoteDetail 获取帖子详情
func (c *XiaoHongShuCrawler) getNoteDetail(noteID, xsecSource, xsecToken string) (*xhs.Note, error) {
	c.logger.Info("[XiaoHongShuCrawler] Getting note detail: %s", noteID)

	note, err := c.client.GetNoteByID(noteID, xsecSource, xsecToken)
	if err != nil {
		// 尝试从HTML获取
		note, err = c.client.GetNoteByIDFromHTML(noteID, xsecSource, xsecToken, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get note detail: %w", err)
		}
	}

	if note == nil {
		return nil, fmt.Errorf("failed to get note detail, note is nil")
	}

	// 添加来源关键词
	note.SourceKeyword = ""

	return note, nil
}

// processCreatorNotes 处理创作者帖子
func (c *XiaoHongShuCrawler) processCreatorNotes(notes []*xhs.Note) error {
	for _, note := range notes {
		if err := c.store.SaveNote(note); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to save creator note: %v", err)
		}
	}
	return nil
}

// getNoteComments 获取帖子评论
func (c *XiaoHongShuCrawler) getNoteComments(noteID, xsecToken string) error {
	c.logger.Info("[XiaoHongShuCrawler] Getting comments for note: %s", noteID)

	return c.client.GetNoteAllComments(
		noteID,
		xsecToken,
		c.config.CrawlerMaxSleepSec,
		c.processComments,
		c.config.CrawlerMaxCommentsCountSingleNotes,
	)
}

// processComments 处理评论
func (c *XiaoHongShuCrawler) processComments(comments []*xhs.Comment) error {
	for _, comment := range comments {
		if err := c.store.SaveComment(comment); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to save comment: %v", err)
		}
	}
	return nil
}

// batchGetNoteComments 批量获取帖子评论
func (c *XiaoHongShuCrawler) batchGetNoteComments(notes []*xhs.Note) {
	// 使用并发获取评论
	semaphore := make(chan struct{}, c.config.MaxConcurrencyNum)

	for _, note := range notes {
		c.wg.Add(1)
		go func(note *xhs.Note) {
			defer c.wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := c.getNoteComments(note.NoteID, note.XsecToken); err != nil {
				c.logger.Error("[XiaoHongShuCrawler] Failed to get comments for note %s: %v", note.NoteID, err)
			}
		}(note)
	}

	c.wg.Wait()
}

// Close 关闭爬虫
func (c *XiaoHongShuCrawler) Close() error {
	c.logger.Info("[XiaoHongShuCrawler] Closing crawler")

	// 关闭浏览器
	if c.browserMgr != nil {
		if err := c.browserMgr.Close(); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to close browser: %v", err)
		}
	}

	// 关闭存储
	if c.store != nil {
		if err := c.store.Close(); err != nil {
			c.logger.Error("[XiaoHongShuCrawler] Failed to close store: %v", err)
		}
	}

	return nil
}

// parseNoteInfoFromNoteURL 解析帖子URL
func parseNoteInfoFromNoteURL(url string) *xhs.NoteURLInfo {
	// 这里实现URL解析逻辑
	return &xhs.NoteURLInfo{
		NoteID:     "test_note_id",
		XsecToken:  "test_token",
		XsecSource: "test_source",
	}
}

// parseCreatorInfoFromURL 解析创作者URL
func parseCreatorInfoFromURL(url string) *xhs.CreatorURLInfo {
	// 这里实现URL解析逻辑
	return &xhs.CreatorURLInfo{
		UserID:     "test_user_id",
		XsecToken:  "test_token",
		XsecSource: "test_source",
	}
}
