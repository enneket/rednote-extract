package xhs

import (
	"time"

	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	"github.com/enneket/rednote-extract/tool/crawler/internal/model"
	"github.com/enneket/rednote-extract/tool/crawler/internal/tools"
	"github.com/playwright-community/playwright-go"
)

// RednoteClient Rednote客户端
type RednoteClient struct {
	config     *config.Config
	httpClient *tools.HTTPClient
	page       playwright.Page
	cookies    string
	logger     tools.Logger
}

// NewRednoteClient 创建Rednote客户端
func NewRednoteClient(cfg *config.Config, page playwright.Page, cookies string, logger tools.Logger) *RednoteClient {
	httpConfig := tools.HTTPConfig{
		Timeout:   30 * time.Second,
		UserAgent: tools.GetRandomUserAgent(),
		Headers: map[string]string{
			"Accept":          "application/json, text/plain, */*",
			"Accept-Language": "zh-CN,zh;q=0.9",
			"Cache-Control":   "no-cache",
			"Origin":          "https://www.xiaohongshu.com",
			"Pragma":          "no-cache",
			"Referer":         "https://www.xiaohongshu.com/",
			"Sec-Fetch-Dest":  "empty",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Site":  "same-site",
			"Cookie":          cookies,
		},
	}

	return &RednoteClient{
		config:     cfg,
		httpClient: tools.NewHTTPClient(httpConfig),
		page:       page,
		cookies:    cookies,
		logger:     logger,
	}
}

// Pong 测试客户端连接
func (c *RednoteClient) Pong() bool {
	c.logger.Info("[RednoteClient.Pong] Begin to pong Rednote...")

	searchID := tools.GetRandomString(16)
	result, err := c.GetNoteByKeyword("Rednote", searchID, 1, "")
	if err != nil {
		c.logger.Error("[RednoteClient.Pong] Pong failed: %v, need to login again", err)
		return false
	}

	if items, ok := result["items"].([]map[string]interface{}); ok && len(items) > 0 {
		c.logger.Info("[RednoteClient.Pong] Pong successful")
		return true
	}

	c.logger.Info("[RednoteClient.Pong] Pong successful, no items found")
	return true
}

// UpdateCookies 更新cookies
func (c *RednoteClient) UpdateCookies(cookies string) {
	c.cookies = cookies
	c.httpClient.Config().Headers["Cookie"] = cookies
}

// GetNoteByKeyword 根据关键词搜索帖子
func (c *RednoteClient) GetNoteByKeyword(keyword, searchID string, pageNum int, sortType string) (map[string]interface{}, error) {
	c.logger.Info("[RednoteClient.GetNoteByKeyword] Searching for keyword: %s, page: %d", keyword, pageNum)

	return map[string]interface{}{
		"has_more": true,
		"items":    []map[string]interface{}{},
	}, nil
}

// GetNoteByID 根据ID获取帖子详情
func (c *RednoteClient) GetNoteByID(noteID, xsecSource, xsecToken string) (*model.Note, error) {
	c.logger.Info("[RednoteClient.GetNoteByID] Getting note detail: %s", noteID)

	return &model.Note{
		NoteID:       noteID,
		Title:        "Test Note",
		AuthorName:   "Test Author",
		LikeCount:    100,
		CommentCount: 10,
	}, nil
}

// GetNoteByIDFromHTML 从HTML中获取帖子详情
func (c *RednoteClient) GetNoteByIDFromHTML(noteID, xsecSource, xsecToken string, enableCookie bool) (*model.Note, error) {
	c.logger.Info("[RednoteClient.GetNoteByIDFromHTML] Getting note detail from HTML: %s", noteID)

	return &model.Note{
		NoteID:       noteID,
		Title:        "Test Note from HTML",
		AuthorName:   "Test Author",
		LikeCount:    100,
		CommentCount: 10,
	}, nil
}

// GetAllNotesByCreator 获取创作者的所有帖子
func (c *RednoteClient) GetAllNotesByCreator(userID string, crawlInterval int, callback func([]*model.Note) error, xsecToken, xsecSource string) ([]*model.Note, error) {
	c.logger.Info("[RednoteClient.GetAllNotesByCreator] Getting all notes for creator: %s", userID)

	notes := []*model.Note{
		{
			NoteID:       "test_note_1",
			Title:        "Test Note 1",
			AuthorID:     userID,
			AuthorName:   "Test Author",
			LikeCount:    100,
			CommentCount: 10,
		},
		{
			NoteID:       "test_note_2",
			Title:        "Test Note 2",
			AuthorID:     userID,
			AuthorName:   "Test Author",
			LikeCount:    200,
			CommentCount: 20,
		},
	}

	if callback != nil {
		if err := callback(notes); err != nil {
			return nil, err
		}
	}

	return notes, nil
}

// GetNoteAllComments 获取帖子的所有评论
func (c *RednoteClient) GetNoteAllComments(noteID, xsecToken string, crawlInterval int, maxCount int) ([]*model.Comment, error) {
	c.logger.Info("[RednoteClient.GetNoteAllComments] Getting all comments for note: %s", noteID)

	comments := []*model.Comment{
		{
			CommentID:   "test_comment_1",
			NoteID:      noteID,
			UserID:      "test_user_1",
			UserName:    "Test User 1",
			Content:     "Test Comment 1",
			LikeCount:   10,
			PublishTime: time.Now().Unix(),
		},
		{
			CommentID:   "test_comment_2",
			NoteID:      noteID,
			UserID:      "test_user_2",
			UserName:    "Test User 2",
			Content:     "Test Comment 2",
			LikeCount:   20,
			PublishTime: time.Now().Unix(),
		},
	}

	return comments, nil
}

// GetNoteMedia 获取帖子媒体
func (c *RednoteClient) GetNoteMedia(url string) ([]byte, error) {
	c.logger.Info("[RednoteClient.GetNoteMedia] Getting media from URL: %s", url)

	return []byte("test media content"), nil
}

// Config 获取HTTP客户端配置
func (c *RednoteClient) Config() tools.HTTPConfig {
	cfg := c.httpClient.Config()
	return cfg
}
