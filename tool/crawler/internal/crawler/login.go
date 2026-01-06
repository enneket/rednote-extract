package xhs

import (
	"encoding/json"

	"github.com/enneket/rednote-extract/tool/crawler/internal/tools"
	"github.com/playwright-community/playwright-go"
)

// RednoteLogin Rednote登录
type RednoteLogin struct {
	cookies    string
	browserCtx playwright.BrowserContext
	logger     tools.Logger
}

// NewRednoteLogin 创建Rednote登录实例
func NewRednoteLogin(cookies string, browserCtx playwright.BrowserContext, logger tools.Logger) *RednoteLogin {
	return &RednoteLogin{
		cookies:    cookies,
		browserCtx: browserCtx,
		logger:     logger,
	}
}

// Begin 开始登录
func (l *RednoteLogin) Begin() error {
	l.logger.Info("[RednoteLogin] Begin login by cookie")

	return l.cookieLogin()
}

// cookieLogin Cookie登录
func (l *RednoteLogin) cookieLogin() error {
	l.logger.Info("[RednoteLogin] Using cookie login")

	if l.cookies == "" {
		l.logger.Info("[RednoteLogin] No cookies provided")
		return nil
	}

	var cookies []playwright.OptionalCookie
	if err := json.Unmarshal([]byte(l.cookies), &cookies); err != nil {
		l.logger.Error("[RednoteLogin] Failed to parse cookies: %v", err)
		return err
	}

	if len(cookies) == 0 {
		l.logger.Info("[RednoteLogin] Empty cookies")
		return nil
	}

	l.logger.Info("[RednoteLogin] Adding %d cookies", len(cookies))

	if err := l.browserCtx.AddCookies(cookies); err != nil {
		l.logger.Error("[RednoteLogin] Failed to add cookies: %v", err)
		return err
	}

	l.logger.Info("[RednoteLogin] Cookie login successful")
	return nil
}
