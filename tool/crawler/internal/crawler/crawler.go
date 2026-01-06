package xhs

import (
	"fmt"
	"net/url"
	"strings"
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
func (c *RednoteCrawler) Init() error {
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

	c.client = NewRednoteClient(c.config, page, c.config.Cookies, c.logger)

	return nil
}

// Search 搜索帖子
func (c *RednoteCrawler) Search() ([]*model.Note, error) {
	c.logger.Info("[RednoteCrawler] Begin search Rednote keywords")

	keyword := c.config.Keyword
	c.logger.Info("[RednoteCrawler] Current search keyword: %s", keyword)

	var allNotes []*model.Note

	page := 1
	searchID := tools.GetRandomString(16)

	for (page-c.config.StartPage+1)*20 <= c.config.CrawlerMaxNotesCount {
		if page < c.config.StartPage {
			c.logger.Info("[RednoteCrawler] Skip page %d", page)
			page++
			continue
		}

		c.logger.Info("[RednoteCrawler] Searching keyword: %s, page: %d", keyword, page)

		// Map config.SortType to SearchSortType
		var sortType SearchSortType
		switch c.config.SortType {
		case "latest":
			sortType = SearchSortTypeLatest
		case "most_liked":
			sortType = SearchSortTypeMostLiked
		default:
			sortType = SearchSortTypeGeneral
		}

		notesRes, err := c.client.GetNoteByKeyword(keyword, searchID, page, 20, sortType, SearchNoteTypeAll)
		if err != nil {
			c.logger.Error("[RednoteCrawler] Failed to get search results: %v", err)
			break
		}

		if items, ok := notesRes["items"].([]map[string]interface{}); ok {
			for _, item := range items {
				note := c.parseNoteFromSearchResult(item, keyword)
				if note != nil {
					allNotes = append(allNotes, note)
				}
			}
		}

		hasMore, _ := notesRes["has_more"].(bool)
		if !hasMore {
			break
		}

		page++

		time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)
	}

	c.logger.Info("[RednoteCrawler] Search completed, found %d notes", len(allNotes))
	return allNotes, nil
}

// parseNoteFromSearchResult 从搜索结果解析帖子
func (c *RednoteCrawler) parseNoteFromSearchResult(item map[string]interface{}, keyword string) *model.Note {
	noteID, _ := item["note_id"].(string)
	if noteID == "" {
		noteID, _ = item["noteId"].(string)
	}
	if noteID == "" {
		return nil
	}

	title, _ := item["title"].(string)
	if title == "" {
		title = "Untitled"
	}

	authorName := ""
	if author, ok := item["user"].(map[string]interface{}); ok {
		authorName, _ = author["nickname"].(string)
	}

	likes := 0
	if likeInfo, ok := item["interact_info"].(map[string]interface{}); ok {
		if likeCount, ok := likeInfo["liked_count"].(float64); ok {
			likes = int(likeCount)
		}
	}

	notes := &model.Note{
		NoteID:     noteID,
		Title:      title,
		AuthorName: authorName,
		LikeCount:  likes,
		SourceType: "search",
		Keyword:    keyword,
	}

	return notes
}

// GetSpecifiedNotes 获取指定帖子详情
func (c *RednoteCrawler) GetSpecifiedNotes() (*model.Note, error) {
	c.logger.Info("[RednoteCrawler] Begin get specified notes")

	noteURL := c.config.Keyword
	c.logger.Info("[RednoteCrawler] Processing note URL: %s", noteURL)

	noteInfo := parseNoteInfoFromNoteURL(noteURL)
	if noteInfo == nil {
		c.logger.Error("[RednoteCrawler] Failed to parse note URL: %s", noteURL)
		return nil, nil
	}

	note, err := c.getNoteDetail(noteInfo.NoteID, noteInfo.XsecSource, noteInfo.XsecToken)
	if err != nil {
		c.logger.Error("[RednoteCrawler] Failed to get note detail: %v", err)
		return nil, err
	}

	note.SourceType = "specified"
	note.Keyword = noteURL

	if err := c.store.SaveNote(note); err != nil {
		c.logger.Error("[RednoteCrawler] Failed to save note: %v", err)
	}

	if c.config.EnableGetComments {
		if comments, err := c.GetNoteComments(note.NoteID, note.XsecToken); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to get comments: %v", err)
		} else {
			for _, comment := range comments {
				if err := c.store.SaveComment(comment); err != nil {
					c.logger.Error("[RednoteCrawler] Failed to save comment: %v", err)
				}
			}
			note.Comments = comments
		}
	}

	time.Sleep(time.Duration(c.config.CrawlerMaxSleepSec) * time.Second)

	return note, nil
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

// GetNoteComments 获取帖子评论
func (c *RednoteCrawler) GetNoteComments(noteID, xsecToken string) ([]*model.Comment, error) {
	c.logger.Info("[RednoteCrawler] Getting comments for note: %s", noteID)

	comments, err := c.client.GetNoteAllComments(
		noteID,
		xsecToken,
		c.config.CrawlerMaxSleepSec,
		c.config.CrawlerMaxCommentsCountSingleNotes,
	)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		if err := c.store.SaveComment(comment); err != nil {
			c.logger.Error("[RednoteCrawler] Failed to save comment: %v", err)
		}
	}

	return comments, nil
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

// parseNoteInfoFromNoteURL 解析帖子URL
func parseNoteInfoFromNoteURL(urlStr string) *model.NoteURLInfo {
	// Parse URL to extract path and query
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return &model.NoteURLInfo{}
	}

	// Extract note ID from path
	noteID := parsedURL.Path
	// Get the last segment of the path
	if idx := strings.LastIndex(noteID, "/"); idx != -1 {
		noteID = noteID[idx+1:]
	}

	// Parse query parameters
	params := make(map[string]string)
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Get xsec_token and xsec_source from params
	xsecToken := params["xsec_token"]
	xsecSource := params["xsec_source"]

	return &model.NoteURLInfo{
		NoteID:     noteID,
		XsecToken:  xsecToken,
		XsecSource: xsecSource,
	}
}
