package tools

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// OrderedMap 保持键值对的插入顺序，用于模拟 Python 3.7+ 的字典行为
type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

// NewOrderedMap 创建一个新的有序 Map
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string]interface{}),
	}
}

// Set 设置键值对
func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Get 获取键对应的值
func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := om.values[key]
	return val, exists
}

// Keys 返回所有键（按插入顺序）
func (om *OrderedMap) Keys() []string {
	return om.keys
}

// Len 返回键值对数量
func (om *OrderedMap) Len() int {
	return len(om.keys)
}

// MarshalJSON 实现 JSON 序列化，保持键的顺序
func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.WriteString("{")
	for i, key := range om.keys {
		if i > 0 {
			buf.WriteString(",")
		}
		// 序列化键
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteString(":")
		// 序列化值
		valueBytes, err := json.Marshal(om.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}
	buf.WriteString("}")
	return []byte(buf.String()), nil
}

func buildSignString(uri string, data interface{}, method string) string {
	if strings.ToUpper(method) == "POST" {
		// POST request uses JSON format
		c := uri
		if data != nil {
			switch v := data.(type) {
			case *OrderedMap:
				// 使用OrderedMap的MarshalJSON保持键的顺序
				jsonBytes, err := json.Marshal(v)
				if err == nil {
					c += string(jsonBytes)
				}
			case map[string]interface{}:
				// 使用与Python json.dumps一致的序列化选项: separators=(",", ":"), ensure_ascii=False
				jsonBytes, err := json.Marshal(v)
				if err == nil {
					c += string(jsonBytes)
				}
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
		case *OrderedMap:
			if v.Len() == 0 {
				return uri
			}
			var params []string
			for _, key := range v.Keys() {
				value, _ := v.Get(key)
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
				// URL encoding - Python的quote(value_str, safe='')
				valueStr = url.QueryEscape(valueStr)
				// url.QueryEscape将空格编码为+，我们需要将其替换为%20
				valueStr = strings.ReplaceAll(valueStr, "+", "%20")
				params = append(params, key+"="+valueStr)
			}
			return uri + "?" + strings.Join(params, "&")
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
				// url.QueryEscape将空格编码为+，我们需要将其替换为%20
				valueStr = strings.ReplaceAll(valueStr, "+", "%20")
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
	// Create the payload struct
	type payload struct {
		X0 string `json:"x0"`
		X1 string `json:"x1"`
		X2 string `json:"x2"`
		X3 string `json:"x3"`
		X4 string `json:"x4"`
	}
	s := payload{
		X0: "4.2.1",
		X1: "xhs-pc-web",
		X2: "Mac OS",
		X3: x3Value,
		X4: dataType,
	}
	// Marshal to JSON
	jsonBytes, _ := json.Marshal(s)
	jsonStr := string(jsonBytes)

	// Escape non-ASCII characters to match Python's json.dumps default behavior
	var builder strings.Builder
	for _, r := range jsonStr {
		if r < 128 {
			builder.WriteRune(r)
		} else {
			builder.WriteString(fmt.Sprintf("\\u%04x", r))
		}
	}
	escapedJSON := builder.String()

	return "XYS_" + base64.StdEncoding.EncodeToString([]byte(escapedJSON))
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

func PreHeadersWithPlaywright(page playwright.Page, urlStr string, cookieDict map[string]string, params *OrderedMap, payload *OrderedMap) (map[string]string, error) {
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
