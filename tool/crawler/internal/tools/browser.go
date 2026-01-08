package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	"github.com/playwright-community/playwright-go"
)

// BrowserManager 浏览器管理器
type BrowserManager struct {
	config     *config.Config
	browser    playwright.Browser
	context    playwright.BrowserContext
	page       playwright.Page
	playwright *playwright.Playwright
	logger     Logger
}

// BrowserResources 浏览器资源管理
type BrowserResources struct {
	Playwright *playwright.Playwright
	Browser    playwright.Browser
	Context    playwright.BrowserContext
	Page       playwright.Page
}

// NewBrowserManager 创建浏览器管理器
func NewBrowserManager(cfg *config.Config, logger Logger) *BrowserManager {
	return &BrowserManager{
		config: cfg,
		logger: logger,
	}
}

// Close 关闭所有资源
func (bm *BrowserManager) Close() error {
	if bm.page != nil {
		if err := bm.page.Close(); err != nil {
			bm.logger.Error("关闭页面失败: %v", err)
		}
		bm.page = nil
	}

	if bm.context != nil {
		if err := bm.context.Close(); err != nil {
			bm.logger.Error("关闭上下文失败: %v", err)
		}
		bm.context = nil
	}

	if bm.browser != nil {
		if err := bm.browser.Close(); err != nil {
			bm.logger.Error("关闭浏览器失败: %v", err)
		}
		bm.browser = nil
	}

	if bm.playwright != nil {
		if err := bm.playwright.Stop(); err != nil {
			bm.logger.Error("关闭Playwright失败: %v", err)
		}
		bm.playwright = nil
	}

	return nil
}

// Launch 启动浏览器
func (bm *BrowserManager) Launch() (playwright.Browser, error) {
	return bm.launchStandard()
}

// launchStandard 标准模式启动浏览器
func (bm *BrowserManager) launchStandard() (playwright.Browser, error) {
	bm.logger.Info("[BrowserManager] Launching browser using standard mode")

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}

	bm.playwright = pw

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(bm.config.Headless),
		Args: []string{
			"--disable-web-security",            // 禁用跨域安全检查（解决跨域访问localStorage）
			"--disable-features=PrivacySandbox", // 禁用隐私沙盒，允许访问存储
			"--allow-running-insecure-content",  // 放宽内容安全限制
		},
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}
	bm.browser = browser

	// 创建浏览器上下文
	context, err := bm.CreateContext()
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	bm.context = context

	bm.logger.Info("[BrowserManager] Browser launched successfully with context")
	return browser, nil
}

// CreateContext 创建浏览器上下文
func (bm *BrowserManager) CreateContext() (playwright.BrowserContext, error) {
	if bm.browser == nil {
		return nil, fmt.Errorf("browser not initialized")
	}

	userDataDir, err := bm.SaveUserDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to save user data dir: %w", err)
	}

	ctx, err := bm.playwright.Chromium.LaunchPersistentContext(userDataDir, playwright.BrowserTypeLaunchPersistentContextOptions{
		Permissions: []string{"storage-access"}, // 授予存储权限
		// 禁用Cookie阻止（关键：允许localStorage写入/读取）
		AcceptDownloads: playwright.Bool(true),
		UserAgent:       playwright.String(bm.config.UserAgent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	bm.context = ctx

	return ctx, nil
}

// NewPage 创建新页面
func (bm *BrowserManager) NewPage() (playwright.Page, error) {
	if bm.context == nil {
		return nil, fmt.Errorf("context not initialized")
	}

	page, err := bm.context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	bm.page = page

	err = page.SetViewportSize(1920, 1080)
	if err != nil {
		page.Close()
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	bm.AddStealthScript(page, "./tool/crawler/internal/js/stealth.min.js")

	return page, nil
}

// AddStealthScript 添加 stealth 反检测脚本
func (bm *BrowserManager) AddStealthScript(page playwright.Page, scriptPath string) error {
	bm.logger.Info("[BrowserManager] Adding stealth script: %s", scriptPath)

	if !strings.HasSuffix(scriptPath, ".js") {
		scriptPath = scriptPath + ".js"
	}

	err := page.AddInitScript(playwright.Script{
		Path: playwright.String(scriptPath),
	})
	if err != nil {
		return fmt.Errorf("添加 stealth 脚本失败: %w", err)
	}

	bm.logger.Info("[BrowserManager] 成功添加 stealth 脚本: %s", scriptPath)
	return nil
}

// GetCookies 获取浏览器cookies
func (bm *BrowserManager) GetCookies(domain string) ([]playwright.Cookie, error) {
	if bm.context == nil {
		return nil, fmt.Errorf("context not initialized")
	}

	cookies, err := bm.context.Cookies()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}

	if domain != "" {
		var filteredCookies []playwright.Cookie
		for _, cookie := range cookies {
			if strings.Contains(cookie.Domain, domain) {
				filteredCookies = append(filteredCookies, cookie)
			}
		}
		return filteredCookies, nil
	}

	return cookies, nil
}

// ConvertCookiesToString 将cookies转换为字符串
func (bm *BrowserManager) ConvertCookiesToString(cookies []playwright.Cookie) string {
	var cookieStr strings.Builder
	for _, cookie := range cookies {
		cookieStr.WriteString(fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	}

	result := cookieStr.String()
	if len(result) > 2 {
		result = result[:len(result)-2]
	}

	return result
}

// AddCookies 添加cookies到上下文
func (bm *BrowserManager) AddCookies(cookies []playwright.OptionalCookie) error {
	if bm.context == nil {
		return fmt.Errorf("context not initialized")
	}

	return bm.context.AddCookies(cookies)
}

// SaveUserData 保存用户数据目录
func (bm *BrowserManager) SaveUserDataDir() (string, error) {
	userDataDir := filepath.Join(".", "browser_data", fmt.Sprintf(bm.config.UserDataDir, "rednote"))
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user data directory: %w", err)
	}
	absUserDataDir, err := filepath.Abs(userDataDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for user data directory: %w", err)
	}
	return absUserDataDir, nil
}

// Navigate 导航到指定URL
func (bm *BrowserManager) Navigate(page playwright.Page, url string) error {
	if page == nil {
		return fmt.Errorf("page not initialized")
	}

	_, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	return nil
}

// WaitForNavigation 等待页面导航完成
func (bm *BrowserManager) WaitForNavigation(page playwright.Page, timeout time.Duration) error {
	if page == nil {
		return fmt.Errorf("page not initialized")
	}

	done := make(chan bool)
	go func() {
		if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State: playwright.LoadStateLoad,
		}); err != nil {
			bm.logger.Error("等待页面加载状态失败: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("navigation timeout")
	}
}

// GetPage 获取当前页面
func (bm *BrowserManager) GetPage() playwright.Page {
	return bm.page
}

// GetContext 获取当前上下文
func (bm *BrowserManager) GetContext() playwright.BrowserContext {
	return bm.context
}

// GetBrowser 获取当前浏览器
func (bm *BrowserManager) GetBrowser() playwright.Browser {
	return bm.browser
}
