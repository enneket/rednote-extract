package xhs

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	"github.com/enneket/rednote-extract/tool/crawler/internal/model"
	"github.com/enneket/rednote-extract/tool/crawler/internal/store"
	"github.com/enneket/rednote-extract/tool/crawler/internal/tools"
	"github.com/playwright-community/playwright-go"
)

// RednoteCrawler Rednote爬虫
type RednoteCrawler struct {
	config     *config.Config
	client     *RednoteClient
	browser    playwright.Browser
	page       playwright.Page
	store      store.Store
	browserMgr *tools.BrowserManager
	logger     tools.Logger
	wg         sync.WaitGroup
}

// NewRednoteCrawler 创建Rednote爬虫
func NewRednoteCrawler(cfg *config.Config, store store.Store, logger tools.Logger) *RednoteCrawler {
	return &RednoteCrawler{
		config: cfg,
		store:  store,
		logger: logger,
	}
}

// Start 启动爬虫
func (c *RednoteCrawler) Start() error {
	c.logger.Info("[RednoteCrawler] Starting crawler")

	c.browserMgr = tools.NewBrowserManager(c.config, c.logger)

	browser, err := c.browserMgr.Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}
	c.browser = browser

	page, err := c.browserMgr.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	c.page = page

	if _, err := c.page.Goto("https://www.xiaohongshu.com"); err != nil {
		return fmt.Errorf("failed to navigate to xiaohongshu: %w", err)
	}

	cookies, err := c.browserMgr.GetCookies(".xiaohongshu.com")
	if err != nil {
		c.logger.Error("[RednoteCrawler] Failed to get cookies: %v", err)
	}
	cookieJSON, err := json.Marshal(cookies)
	if err != nil {
		c.logger.Error("[RednoteCrawler] Failed to marshal cookies: %v", err)
	}

	c.client = NewRednoteClient(c.config, page, string(cookieJSON), c.logger)

	if !c.client.Pong() {
		if err := c.login(); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}
	}

	switch c.config.CrawlerType {
	case "search":
		return c.search()
	case "detail":
		return c.getSpecifiedNotes()
	default:
		return fmt.Errorf("unsupported crawler type: %s", c.config.CrawlerType)
	}
}

// login 登录Rednote
func (c *RednoteCrawler) login() error {
	c.logger.Info("[RednoteCrawler] Login required")

	ctx := c.page.Context()

	cookies, err := c.browserMgr.GetCookies(".xiaohongshu.com")
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}
	cookieJSON, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	loginObj := NewRednoteLogin(string(cookieJSON), ctx, c.logger)
	if err := loginObj.Begin(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	newCookies, err := c.browserMgr.GetCookies(".xiaohongshu.com")
	if err != nil {
		return fmt.Errorf("failed to get cookies after login: %w", err)
	}
	newCookieJSON, err := json.Marshal(newCookies)
	if err != nil {
		return fmt.Errorf("failed to marshal new cookies: %w", err)
	}
	c.client.UpdateCookies(string(newCookieJSON))

	return nil
}

// search 搜索帖子
func (c *RednoteCrawler) search() error {
	c.logger.Info("[RednoteCrawler] Begin search Rednote keywords")

	for _, keyword := range c.config.Keywords {
		c.logger.Info("[RednoteCrawler] Current search keyword: %s", keyword)

		page := 1
		searchID := tools.GetRandomString(16)

		for (page-c.config.StartPage+1)*20 <= c.config.CrawlerMaxNotesCount {
			if page < c.config.StartPage {
				c.logger.Info("[RednoteCrawler] Skip page %d", page)
				page++
				continue
			}

			c.logger.Info("[RednoteCrawler] Searching keyword: %s, page: %d", keyword, page)

			notesRes, err := c.client.GetNoteByKeyword(keyword, searchID, page, c.config.SortType)
			if err != nil {
				c.logger.Error("[RednoteCrawler] Failed to get search results: %v", err)
				break
			}

			_ = notesRes

			page++

			time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
		}
	}

	return nil
}

// getSpecifiedNotes 获取指定帖子详情
func (c *RednoteCrawler) getSpecifiedNotes() error {
	c.logger.Info("[RednoteCrawler] Begin get specified notes")

	for _, noteURL := range c.config.Keywords {
		c.logger.Info("[RednoteCrawler] Processing note URL: %s", noteURL)

		noteInfo := parseNoteInfoFromNoteURL(noteURL)
		if noteInfo == nil {
			c.logger.Error("[RednoteCrawler] Failed to parse note URL: %s", noteURL)
			continue
		}

		note, err := c.getNoteDetail(noteInfo.NoteID, noteInfo.XsecSource, noteInfo.XsecToken)
		if err != nil {
			c.logger.Error("[RednoteCrawler] Failed to get note detail: %v", err)
			continue
		}

		if err := c.store.SaveNote(note); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to save note: %v", err)
		}

		if c.config.EnableGetComments {
			if err := c.getNoteComments(note.NoteID, note.XsecToken); err != nil {
				c.logger.Error("[RednoteCrawler] Failed to get comments: %v", err)
			}
		}

		time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
	}

	return nil
}

// getNoteDetail 获取帖子详情
func (c *RednoteCrawler) getNoteDetail(noteID, xsecSource, xsecToken string) (*model.Note, error) {
	c.logger.Info("[RednoteCrawler] Getting note detail: %s", noteID)

	note, err := c.client.GetNoteByID(noteID, xsecSource, xsecToken)
	if err != nil {
		note, err = c.client.GetNoteByIDFromHTML(noteID, xsecSource, xsecToken, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get note detail: %w", err)
		}
	}

	if note == nil {
		return nil, fmt.Errorf("failed to get note detail, note is nil")
	}

	note.SourceKeyword = ""

	return note, nil
}

// processCreatorNotes 处理创作者帖子
func (c *RednoteCrawler) processCreatorNotes(notes []*model.Note) error {
	for _, note := range notes {
		if err := c.store.SaveNote(note); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to save creator note: %v", err)
		}
	}
	return nil
}

// getNoteComments 获取帖子评论
func (c *RednoteCrawler) getNoteComments(noteID, xsecToken string) error {
	c.logger.Info("[RednoteCrawler] Getting comments for note: %s", noteID)

	return c.client.GetNoteAllComments(
		noteID,
		xsecToken,
		c.config.CrawlerMaxSleepSec,
		c.processComments,
		c.config.CrawlerMaxCommentsCountSingleNotes,
	)
}

// processComments 处理评论
func (c *RednoteCrawler) processComments(comments []*model.Comment) error {
	for _, comment := range comments {
		if err := c.store.SaveComment(comment); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to save comment: %v", err)
		}
	}
	return nil
}

// batchGetNoteComments 批量获取帖子评论
func (c *RednoteCrawler) batchGetNoteComments(notes []*model.Note) {
	semaphore := make(chan struct{}, c.config.MaxConcurrencyNum)

	for _, note := range notes {
		c.wg.Add(1)
		go func(note *model.Note) {
			defer c.wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := c.getNoteComments(note.NoteID, note.XsecToken); err != nil {
				c.logger.Error("[RednoteCrawler] Failed to get comments for note %s: %v", note.NoteID, err)
			}
		}(note)
	}

	c.wg.Wait()
}

// Close 关闭爬虫
func (c *RednoteCrawler) Close() error {
	c.logger.Info("[RednoteCrawler] Closing crawler")

	if c.browserMgr != nil {
		if err := c.browserMgr.Close(); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to close browser: %v", err)
		}
	}

	if c.store != nil {
		if err := c.store.Close(); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to close store: %v", err)
		}
	}

	return nil
}

// NoteURLInfo 帖子URL信息
type NoteURLInfo struct {
	NoteID     string
	XsecToken  string
	XsecSource string
}

// CreatorURLInfo 创作者URL信息
type CreatorURLInfo struct {
	UserID     string
	XsecToken  string
	XsecSource string
}

// parseNoteInfoFromNoteURL 解析帖子URL
func parseNoteInfoFromNoteURL(url string) *NoteURLInfo {
	return &NoteURLInfo{
		NoteID:     "test_note_id",
		XsecToken:  "test_token",
		XsecSource: "test_source",
	}
}

// parseCreatorInfoFromURL 解析创作者URL
func parseCreatorInfoFromURL(url string) *CreatorURLInfo {
	return &CreatorURLInfo{
		UserID:     "test_user_id",
		XsecToken:  "test_token",
		XsecSource: "test_source",
	}
}
