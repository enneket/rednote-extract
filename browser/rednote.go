package browser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/enneket/rednote-extract/utils/gptr"
)

type RednoteNote struct {
	Url      string   // 笔记的url
	Title    string   // 笔记的标题
	Content  string   // 笔记的内容
	Comments []string // 笔记的评论
}

func SearchRednote(ctx context.Context, keyword string) ([]*RednoteNote, error) {
	rednoteUrl := fmt.Sprintf("https://www.xiaohongshu.com/search_result/?keyword=%s", keyword)

	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		return nil, err
	}
	// 先跳转到搜索页
	_, err = but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &rednoteUrl,
	})
	if err != nil {
		return nil, err
	}
	// 再获取页面的 HTML
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionExtractContent,
		Goal:   gptr.Of(keyword),
	})
	if err != nil {
		return nil, err
	}
	// 解析 HTML 内容
	fmt.Println(result.Output)

	but.Cleanup()

	// 提取笔记URL和标题
	notes, err := ExtractRednoteUrlAndTitle(ctx, result.Output)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func GetRednoteContentAndComments(note *RednoteNote) error {
	// 打开小红书
	// 输入笔记url
	// 查看笔记内容
	return nil
}

func ExtractRednoteUrlAndTitle(ctx context.Context, htmlContent string) ([]*RednoteNote, error) {
	if strings.TrimSpace(htmlContent) == "" {
		return nil, errors.New("HTML内容不能为空")
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败：%w", err)
	}

	// 定位所有有效笔记项
	noteItemXPath := "//section[contains(@class, 'note-item') and not(.//div[contains(@class, 'query-note-wrapper')])]"
	// 每个笔记项中，定位有效链接
	noteUrlXPath := ".//a[starts-with(@href, '/explore/') and @style='display: none;']/@href"
	// 每个笔记项中，定位笔记标题
	noteTitleXPath := ".//div[contains(@class, 'footer')]//a[contains(@class, 'title')]//span/text()"

	// 提取所有有效笔记项节点
	noteItems, err := htmlquery.QueryAll(doc, noteItemXPath)
	if err != nil {
		return nil, fmt.Errorf("查询笔记项节点失败：%w", err)
	}
	if len(noteItems) == 0 {
		return nil, errors.New("未提取到有效笔记节点")
	}

	// 遍历笔记项，提取URL和标题
	var result []*RednoteNote
	xiaohongshuDomain := "https://www.xiaohongshu.com" // 小红书基础域名，用于补全绝对URL

	for _, item := range noteItems {
		// 提取笔记相对URL
		urlNode := htmlquery.FindOne(item, noteUrlXPath)
		if urlNode == nil {
			continue // 无有效URL，跳过当前笔记项
		}
		relativeURL := strings.TrimSpace(htmlquery.InnerText(urlNode))
		if relativeURL == "" {
			continue
		}

		// 提取笔记标题
		titleNode := htmlquery.FindOne(item, noteTitleXPath)
		title := ""
		if titleNode != nil {
			title = strings.TrimSpace(htmlquery.InnerText(titleNode))
		}

		// 补全绝对URL
		absoluteURL := fmt.Sprintf("%s%s", xiaohongshuDomain, relativeURL)

		// 存入结果切片
		result = append(result, &RednoteNote{
			Url:   absoluteURL,
			Title: title,
		})
	}

	// 返回结果
	return result, nil
}
