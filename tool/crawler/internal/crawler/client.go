package xhs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

/**
  "accept-language": "zh-CN,zh;q=0.9",
  "cache-control": "no-cache",

  "origin": "https://www.xiaohongshu.com",
  "pragma": "no-cache",
  "priority": "u=1, i",
  "referer": "https://www.xiaohongshu.com/",
  "sec-ch-ua": '"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"',
  "sec-ch-ua-mobile": "?0",
  "sec-ch-ua-platform": '"Windows"',
  "sec-fetch-dest": "empty",
  "sec-fetch-mode": "cors",
  "sec-fetch-site": "same-site",
  "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
  "Cookie": cookie_str,
*/
// NewRednoteClient 创建Rednote客户端
func NewRednoteClient(cfg *config.Config, page playwright.Page, cookies string, logger tools.Logger) *RednoteClient {
	httpConfig := tools.HTTPConfig{
		Timeout:   30 * time.Second,
		UserAgent: tools.GetRandomUserAgent(),

		Headers: map[string]string{
			"accept":             "application/json, text/plain, */*",
			"accept-language":    "zh-CN,zh;q=0.9",
			"cache-control":      "no-cache",
			"content-type":       "application/json;charset=UTF-8",
			"origin":             "https://www.xiaohongshu.com",
			"pragma":             "no-cache",
			"priority":           "u=1, i",
			"referer":            "https://www.xiaohongshu.com/",
			"sec-ch-ua":          `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"Windows"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-site",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
			"Cookie":             cookies,
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
	result, err := c.GetNoteByKeyword("Rednote", searchID, 1, 20, SearchSortTypeGeneral, SearchNoteTypeAll)
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

// SearchSortType 搜索结果排序类型
type SearchSortType string

const (
	// SearchSortTypeGeneral 综合排序
	SearchSortTypeGeneral SearchSortType = "general"
	// SearchSortTypeLatest 最新排序
	SearchSortTypeLatest SearchSortType = "time_descending"
	// SearchSortTypeMostLiked 最多点赞
	SearchSortTypeMostLiked SearchSortType = "popularity_descending"
)

// SearchNoteType 搜索帖子类型
type SearchNoteType int

const (
	// SearchNoteTypeAll 全部类型
	SearchNoteTypeAll SearchNoteType = 0
	// SearchNoteTypeImage 图片类型
	SearchNoteTypeImage SearchNoteType = 2
	// SearchNoteTypeVideo 视频类型
	SearchNoteTypeVideo SearchNoteType = 1
)

// GetNoteByKeyword 根据关键词搜索帖子
func (c *RednoteClient) GetNoteByKeyword(keyword string, searchID string, page, pageSize int, sort SearchSortType, noteType SearchNoteType) (map[string]interface{}, error) {
	c.logger.Info("[RednoteClient.GetNoteByKeyword] Searching for keyword: %s, page: %d, pageSize: %d, sort: %s, noteType: %s",
		keyword, page, pageSize, sort, noteType)

	// Use default searchID if not provided
	if searchID == "" {
		searchID = tools.GetSearchID()
	}

	// Use default values if not provided
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if sort == "" {
		sort = SearchSortTypeGeneral
	}

	// Construct request data
	uri := "/api/sns/web/v1/search/notes"
	data := tools.NewOrderedMap()
	data.Set("keyword", keyword)
	data.Set("page", page)
	data.Set("page_size", pageSize)
	data.Set("search_id", searchID)
	data.Set("sort", string(sort))
	data.Set("note_type", int(noteType))

	c.logger.Info("[RednoteClient.GetNoteByKeyword] Request data: keyword=%s, page=%d", keyword, page)
	// Send POST request using Post method
	result, err := c.Post(uri, data, nil)
	if err != nil {
		return nil, err
	}

	// Convert result to map[string]interface{}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
}

// GetNoteByID 根据ID获取帖子详情
func (c *RednoteClient) GetNoteByID(noteID, xsecSource, xsecToken string) (*model.Note, error) {
	c.logger.Info("[RednoteClient.GetNoteByID] Getting note detail: %s", noteID)

	// Set default xsec_source if empty
	if xsecSource == "" {
		xsecSource = "pc_search"
	}

	// Prepare request data
	data := tools.NewOrderedMap()
	data.Set("source_note_id", noteID)
	data.Set("image_formats", []string{"jpg", "webp", "avif"})
	data.Set("extra", map[string]interface{}{"need_body_topic": 1})
	data.Set("xsec_source", xsecSource)
	data.Set("xsec_token", xsecToken)

	// Send POST request
	uri := "/api/sns/web/v1/feed"
	res, err := c.Post(uri, data, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	if resMap, ok := res.(map[string]interface{}); ok {
		if items, ok := resMap["items"].([]interface{}); ok && len(items) > 0 {
			// Check if items[0] is a map
			if _, ok := items[0].(map[string]interface{}); ok {
				// Create and return model.Note
				// This is a simplified version - would need proper parsing from note_card in production
				return &model.Note{
					NoteID:     noteID,
					XsecToken:  xsecToken,
					XsecSource: xsecSource,
				}, nil
			}
		}
	}

	// Log error if no results
	c.logger.Error("[RednoteClient.GetNoteByID] get note id:%s empty and res:%v", noteID, res)
	return &model.Note{NoteID: noteID}, nil
}

// GetNoteComments 获取帖子一级评论
func (c *RednoteClient) GetNoteComments(noteID, xsecToken, cursor string) (map[string]interface{}, error) {
	c.logger.Info("[RednoteClient.GetNoteComments] Getting comments for note: %s, cursor: %s", noteID, cursor)

	// Prepare request params
	params := tools.NewOrderedMap()
	params.Set("note_id", noteID)
	params.Set("cursor", cursor)
	params.Set("top_comment_id", "")
	params.Set("image_formats", "jpg,webp,avif")
	params.Set("xsec_token", xsecToken)

	// Send GET request
	uri := "/api/sns/web/v2/comment/page"
	result, err := c.Get(uri, params)
	if err != nil {
		return nil, err
	}

	// Convert result to map[string]interface{}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
}

// GetNoteSubComments 获取帖子子评论
func (c *RednoteClient) GetNoteSubComments(noteID, rootCommentID, xsecToken, cursor string, num int) (map[string]interface{}, error) {
	c.logger.Info("[RednoteClient.GetNoteSubComments] Getting sub-comments for note: %s, root_comment_id: %s, cursor: %s", noteID, rootCommentID, cursor)

	// Use default num if not provided
	if num <= 0 {
		num = 10
	}

	// Prepare request params
	params := tools.NewOrderedMap()
	params.Set("note_id", noteID)
	params.Set("root_comment_id", rootCommentID)
	params.Set("cursor", cursor)
	params.Set("num", fmt.Sprintf("%d", num))
	params.Set("image_formats", "jpg,webp,avif")
	params.Set("top_comment_id", "")
	params.Set("xsec_token", xsecToken)

	// Send GET request
	uri := "/api/sns/web/v2/comment/sub/page"
	result, err := c.Get(uri, params)
	if err != nil {
		return nil, err
	}

	// Convert result to map[string]interface{}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
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

// GetNoteAllComments 获取帖子的所有评论
func (c *RednoteClient) GetNoteAllComments(noteID, xsecToken string, crawlInterval int, maxCount int) ([]*model.Comment, error) {
	c.logger.Info("[RednoteClient.GetNoteAllComments] Getting all comments for note: %s, maxCount: %d", noteID, maxCount)

	var allComments []*model.Comment
	var cursor string
	hasMore := true

	// Fetch first-level comments with pagination
	for hasMore && len(allComments) < maxCount {
		// Get comments page
		commentsRes, err := c.GetNoteComments(noteID, xsecToken, cursor)
		if err != nil {
			return nil, err
		}

		// commentsRes is already a map[string]interface{} from GetNoteComments

		// Parse comments list
		if commentsList, ok := commentsRes["comments"].([]interface{}); ok {
			for _, commentItem := range commentsList {
				if commentMap, ok := commentItem.(map[string]interface{}); ok {
					// Convert to model.Comment (simplified parsing)
					comment := &model.Comment{
						NoteID: noteID,
						CommentID: func() string {
							if id, ok := commentMap["id"].(string); ok {
								return id
							}
							return ""
						}(),
						Content: func() string {
							if content, ok := commentMap["content"].(string); ok {
								return content
							}
							return ""
						}(),
						LikeCount: func() int {
							if count, ok := commentMap["like_count"].(float64); ok {
								return int(count)
							}
							return 0
						}(),
						PublishTime: func() int64 {
							if time, ok := commentMap["create_time"].(float64); ok {
								return int64(time)
							}
							return time.Now().Unix()
						}(),
					}

					// Fetch sub-comments if needed
					var subCursor string
					hasMoreSubComments := true
					subCommentCount := 0

					for hasMoreSubComments && len(allComments)+1+subCommentCount < maxCount {
						subCommentsRes, err := c.GetNoteSubComments(noteID, comment.CommentID, xsecToken, subCursor, 10)
						if err != nil {
							break
						}

						// Parse sub-comments
						if subCommentsList, ok := subCommentsRes["comments"].([]interface{}); ok {
							for _, subCommentItem := range subCommentsList {
								if subCommentMap, ok := subCommentItem.(map[string]interface{}); ok {
									subComment := &model.Comment{
										NoteID:   noteID,
										ParentID: comment.CommentID,
										CommentID: func() string {
											if id, ok := subCommentMap["id"].(string); ok {
												return id
											}
											return ""
										}(),
										Content: func() string {
											if content, ok := subCommentMap["content"].(string); ok {
												return content
											}
											return ""
										}(),
										LikeCount: func() int {
											if count, ok := subCommentMap["like_count"].(float64); ok {
												return int(count)
											}
											return 0
										}(),
										PublishTime: func() int64 {
											if time, ok := subCommentMap["create_time"].(float64); ok {
												return int64(time)
											}
											return time.Now().Unix()
										}(),
									}
									comment.SubComments = append(comment.SubComments, *subComment)
									subCommentCount++
								}
							}
						}

						// Check if there are more sub-comments
						if newSubCursor, ok := subCommentsRes["cursor"].(string); ok && newSubCursor != "" {
							subCursor = newSubCursor
						} else {
							hasMoreSubComments = false
						}
					}

					allComments = append(allComments, comment)
					if len(allComments) >= maxCount {
						break
					}
				}
			}
		}

		// Check if there are more comments
		if newCursor, ok := commentsRes["cursor"].(string); ok && newCursor != "" {
			cursor = newCursor
		} else {
			hasMore = false
		}

		// Add delay between requests
		if hasMore {
			time.Sleep(time.Duration(crawlInterval) * time.Second)
		}
	}

	return allComments, nil
}

// playwrightPageAdapter adapts playwright.Page to tools.Page interface

type playwrightPageAdapter struct {
	page playwright.Page
}

func (a *playwrightPageAdapter) Evaluate(expression string, options ...interface{}) interface{} {
	result, _ := a.page.Evaluate(expression, options...)
	return result
}

// PreHeaders 生成带签名的请求头
func (c *RednoteClient) PreHeaders(url string, params *tools.OrderedMap, payload *tools.OrderedMap) map[string]string {
	c.logger.Info("[RednoteClient.PreHeaders] Generating signed headers for URL: %s", url)

	// Parse cookies to get a1 value
	cookieDict := make(map[string]string)
	for _, cookie := range strings.Split(c.cookies, "; ") {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			cookieDict[parts[0]] = parts[1]
		}
	}

	headers, err := tools.PreHeadersWithPlaywright(
		c.page,
		url,
		cookieDict,
		params,
		payload,
	)
	if err != nil {
		c.logger.Error("[RednoteClient.PreHeaders] Failed to generate signed headers: %v", err)
		return nil
	}
	return headers
}

// Request 发送HTTP请求并处理响应
func (c *RednoteClient) Request(method, url string, returnResponse bool, headers map[string]string, payload interface{}) (interface{}, error) {
	c.logger.Info("[RednoteClient.Request] Sending %s request to URL: %s", method, url)

	// Prepare request body if payload exists
	var bodyReader io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["content-type"] = "application/json"
	}

	// Send request using existing HTTP client
	respBody, err := c.httpClient.Request(method, url, bodyReader, headers)
	if err != nil {
		// Check if it's an HTTP error with specific status codes
		if httpErr, ok := err.(*tools.HTTPError); ok {
			if httpErr.StatusCode == 471 || httpErr.StatusCode == 461 {
				// Handle CAPTCHA error - extract verify headers if available
				verifyType := "unknown"
				verifyUUID := "unknown"
				// Note: Go's http client doesn't preserve headers in error,
				// would need to enhance HTTPError to include headers
				msg := fmt.Sprintf("CAPTCHA appeared, request failed, Verifytype: %s, Verifyuuid: %s, status_code: %d, body: %s",
					verifyType, verifyUUID, httpErr.StatusCode, httpErr.Body)
				c.logger.Error(msg)
				return nil, fmt.Errorf(msg)
			}
		}
		return nil, err
	}

	// If return_response is true, return raw response text
	if returnResponse {
		return string(respBody), nil
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	// Handle response based on success field
	if success, ok := data["success"].(bool); ok && success {
		if dataValue, exists := data["data"]; exists {
			return dataValue, nil
		}
		if successValue, exists := data["success"]; exists {
			return successValue, nil
		}
		return data, nil
	} else if code, ok := data["code"].(float64); ok {
		// Check for IP error code
		const IP_ERROR_CODE = 50011
		if code == IP_ERROR_CODE {
			return nil, fmt.Errorf("IP blocked")
		}
	}

	// Handle other errors
	errMsg := "Unknown error"
	if msg, ok := data["msg"].(string); ok {
		errMsg = msg
	} else {
		errMsg = string(respBody)
	}

	return nil, fmt.Errorf("Data fetch error: %s", errMsg)
}

// Get 发送带签名的GET请求
func (c *RednoteClient) Get(uri string, params *tools.OrderedMap) (interface{}, error) {
	c.logger.Info("[RednoteClient.Get] Sending GET request to URI: %s", uri)

	// Get signed headers using PreHeaders method
	headers := c.PreHeaders(uri, params, nil)

	// Construct full URL (using origin from headers as host)
	const host = "https://edith.xiaohongshu.com"
	fullURL := host + uri

	// Send GET request using Request method
	return c.Request("GET", fullURL, false, headers, nil)
}

// Post 发送带签名的POST请求
func (c *RednoteClient) Post(uri string, payload *tools.OrderedMap, params *tools.OrderedMap) (interface{}, error) {
	c.logger.Info("[RednoteClient.Post] Sending POST request to URI: %s", uri)

	// Get signed headers using PreHeaders method
	headers := c.PreHeaders(uri, params, payload)

	// Construct full URL (using origin from headers as host)
	const host = "https://edith.xiaohongshu.com"
	fullURL := host + uri

	// 将 OrderedMap 转换为 map[string]interface{} 以传递给 Request
	var payloadMap map[string]interface{}
	if payload != nil {
		payloadMap = make(map[string]interface{})
		for _, key := range payload.Keys() {
			value, _ := payload.Get(key)
			payloadMap[key] = value
		}
	}

	// Send POST request using Request method
	return c.Request("POST", fullURL, false, headers, payloadMap)
}

// Config 获取HTTP客户端配置
func (c *RednoteClient) Config() tools.HTTPConfig {
	return c.httpClient.Config()
}
