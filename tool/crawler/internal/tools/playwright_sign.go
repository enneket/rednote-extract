package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func buildSignString(uri string, data interface{}, method string) string {
	if strings.ToUpper(method) == "POST" {
		// POST request uses JSON format
		c := uri
		if data != nil {
			switch v := data.(type) {
			case map[string]interface{}:
				jsonBytes, _ := json.Marshal(v)
				c += string(jsonBytes)
			case string:
				c += v
			}
		}
		return c
	} else {
		// GET request uses query string format
		if data == nil {
			return uri
		}

		switch v := data.(type) {
		case map[string]interface{}:
			if len(v) == 0 {
				return uri
			}
			var params []string
			for key, value := range v {
				var valueStr string
				switch val := value.(type) {
				case []interface{}:
					var strSlice []string
					for _, item := range val {
						strSlice = append(strSlice, fmt.Sprintf("%v", item))
					}
					valueStr = strings.Join(strSlice, ",")
				case nil:
					valueStr = ""
				default:
					valueStr = fmt.Sprintf("%v", val)
				}
				// URL encoding without safe characters
				valueStr = url.QueryEscape(valueStr)
				params = append(params, key+"="+valueStr)
			}
			return uri + "?" + strings.Join(params, "&")
		case string:
			return uri + "?" + v
		default:
			return uri
		}
	}
}

func md5Hex(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func buildXSPayload(x3Value string, dataType string) string {
	if dataType == "" {
		dataType = "object"
	}
	s := map[string]interface{}{
		"x0": "4.2.1",
		"x1": "xhs-pc-web",
		"x2": "Mac OS",
		"x3": x3Value,
		"x4": dataType,
	}
	jsonBytes, _ := json.Marshal(s)
	return "XYS_" + B64Encode(EncodeUtf8(string(jsonBytes)))
}

func buildXSCommon(a1 string, b1 string, x_s string, x_t string) string {
	payload := map[string]interface{}{
		"s0":  3,
		"s1":  "",
		"x0":  "1",
		"x1":  "4.2.2",
		"x2":  "Mac OS",
		"x3":  "xhs-pc-web",
		"x4":  "4.74.0",
		"x5":  a1,
		"x6":  x_t,
		"x7":  x_s,
		"x8":  b1,
		"x9":  Mrc(x_t + x_s + b1),
		"x10": 154,
		"x11": "normal",
	}
	jsonBytes, _ := json.Marshal(payload)
	return B64Encode(EncodeUtf8(string(jsonBytes)))
}

func GetB1FromLocalStorage(page playwright.Page) (string, error) {
	localStorage, err := page.Evaluate(`() => window.localStorage`, nil)
	if err != nil {
		return "", err
	}

	if localStorageMap, ok := localStorage.(map[string]interface{}); ok {
		if b1, exists := localStorageMap["b1"]; exists {
			if b1Str, ok := b1.(string); ok {
				return b1Str, nil
			}
		} else {
			return "", fmt.Errorf("b1 not found in localStorage")
		}
	} else {
		return "", fmt.Errorf("localStorage is not a map[string]interface{}")
	}
	return "", nil
}

func CallMNSV2(page playwright.Page, signStr string, md5Str string) (string, error) {
	// Escape special characters for JavaScript
	signStrEscaped := strings.ReplaceAll(signStr, "\\", "\\\\")
	signStrEscaped = strings.ReplaceAll(signStrEscaped, "'", "\\'")
	signStrEscaped = strings.ReplaceAll(signStrEscaped, "\n", "\\n")

	md5StrEscaped := strings.ReplaceAll(md5Str, "\\", "\\\\")
	md5StrEscaped = strings.ReplaceAll(md5StrEscaped, "'", "\\'")

	evalStr := fmt.Sprintf("window.mnsv2('%s', '%s')", signStrEscaped, md5StrEscaped)
	result, err := page.Evaluate(evalStr, nil)
	if err != nil {
		return "", err
	}

	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}
	return "", nil
}

func SignXSWithPlaywright(page playwright.Page, uri string, data interface{}, method string) (string, error) {
	signStr := buildSignString(uri, data, method)
	md5Str := md5Hex(signStr)
	x3Value, err := CallMNSV2(page, signStr, md5Str)
	if err != nil {
		return "", err
	}

	dataType := "object"
	switch data.(type) {
	case string:
		dataType = "string"
	}

	return buildXSPayload(x3Value, dataType), nil
}

func SignWithPlaywright(page playwright.Page, uri string, data interface{}, a1 string, method string) (map[string]string, error) {
	b1, err := GetB1FromLocalStorage(page)
	if err != nil {
		return nil, err
	}

	x_s, err := SignXSWithPlaywright(page, uri, data, method)
	if err != nil {
		return nil, err
	}

	x_t := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return map[string]string{
		"x-s":          x_s,
		"x-t":          x_t,
		"x-s-common":   buildXSCommon(a1, b1, x_s, x_t),
		"x-b3-traceid": GetTraceId(),
	}, nil
}

func PreHeadersWithPlaywright(page playwright.Page, urlStr string, cookieDict map[string]string, params map[string]interface{}, payload map[string]interface{}) (map[string]string, error) {
	a1Value := cookieDict["a1"]
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	uri := parsedURL.Path

	var data interface{}
	var method string

	if params != nil {
		data = params
		method = "GET"
	} else if payload != nil {
		data = payload
		method = "POST"
	} else {
		return nil, fmt.Errorf("params or payload is required")
	}

	page.Goto("https://www.xiaohongshu.com", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})

	time.Sleep(5 * time.Second)

	signs, err := SignWithPlaywright(page, uri, data, a1Value, method)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"X-S":          signs["x-s"],
		"X-T":          signs["x-t"],
		"x-S-Common":   signs["x-s-common"],
		"X-B3-Traceid": signs["x-b3-traceid"],
	}, nil
}
