package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Page interface {
	Evaluate(expression string, options ...interface{}) interface{}
}

type RequestData interface{}

type PlaywrightClient struct {
	page Page
}

func NewPlaywrightClient(page Page) *PlaywrightClient {
	return &PlaywrightClient{page: page}
}

func BuildSignString(uri string, data interface{}, method string) string {
	methodUpper := strings.ToUpper(method)

	if methodUpper == "POST" {
		c := uri
		if data != nil {
			switch v := data.(type) {
			case map[string]interface{}:
				c += MapToJson(v)
			case string:
				c += v
			case []byte:
				c += string(v)
			}
		}
		return c
	} else {
		if data == nil {
			return uri
		}

		switch v := data.(type) {
		case map[string]string:
			if len(v) == 0 {
				return uri
			}
			params := []string{}
			keys := make([]string, 0, len(v))
			for key := range v {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				valueStr := v[key]
				valueStr = CustomQuote(valueStr, false)
				params = append(params, fmt.Sprintf("%s=%s", key, valueStr))
			}
			return uri + "?" + strings.Join(params, "&")
		case string:
			return uri + "?" + v
		}
		return uri
	}
}

func MapToJson(m map[string]interface{}) string {
	result, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(result)
}

func Md5Hex(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func BuildXsPayload(x3Value string, dataType string) string {
	s := map[string]interface{}{
		"x0": "4.2.1",
		"x1": "xhs-pc-web",
		"x2": "Mac OS",
		"x3": x3Value,
		"x4": dataType,
	}
	jsonStr, _ := json.Marshal(s)
	return "XYS_" + B64Encode(EncodeUtf8(string(jsonStr)))
}

func BuildXsCommon(a1 string, b1 string, xS string, xT string) string {
	payload := map[string]interface{}{
		"s0":  3,
		"s1":  "",
		"x0":  "1",
		"x1":  "4.2.2",
		"x2":  "Mac OS",
		"x3":  "xhs-pc-web",
		"x4":  "4.74.0",
		"x5":  a1,
		"x6":  xT,
		"x7":  xS,
		"x8":  b1,
		"x9":  MRC(xT + xS + b1),
		"x10": 154,
		"x11": "normal",
	}
	jsonStr, _ := json.Marshal(payload)
	return B64Encode(EncodeUtf8(string(jsonStr)))
}

func (c *PlaywrightClient) GetB1FromLocalStorage() string {
	if c.page == nil {
		return ""
	}
	result := c.page.Evaluate("() => window.localStorage")
	if result == nil {
		return ""
	}

	if m, ok := result.(map[string]interface{}); ok {
		if b1, ok := m["b1"]; ok {
			if b1Str, ok := b1.(string); ok {
				return b1Str
			}
		}
	}
	return ""
}

func (c *PlaywrightClient) CallMnsv2(signStr string, md5Str string) string {
	if c.page == nil {
		return ""
	}

	signStrEscaped := strings.ReplaceAll(signStr, "\\", "\\\\")
	signStrEscaped = strings.ReplaceAll(signStrEscaped, "'", "\\'")
	signStrEscaped = strings.ReplaceAll(signStrEscaped, "\n", "\\n")
	md5StrEscaped := strings.ReplaceAll(md5Str, "\\", "\\\\")
	md5StrEscaped = strings.ReplaceAll(md5StrEscaped, "'", "\\'")

	script := fmt.Sprintf("window.mnsv2('%s', '%s')", signStrEscaped, md5StrEscaped)
	result := c.page.Evaluate(script)
	if result == nil {
		return ""
	}
	if resultStr, ok := result.(string); ok {
		return resultStr
	}
	return ""
}

func (c *PlaywrightClient) SignXsWithPlaywright(uri string, data interface{}, method string) string {
	signStr := BuildSignString(uri, data, method)
	md5Str := Md5Hex(signStr)
	x3Value := c.CallMnsv2(signStr, md5Str)

	var dataType string
	switch data.(type) {
	case map[string]interface{}, []interface{}:
		dataType = "object"
	default:
		dataType = "string"
	}

	return BuildXsPayload(x3Value, dataType)
}

func (c *PlaywrightClient) SignWithPlaywright(uri string, data interface{}, a1 string, method string) map[string]interface{} {
	b1 := c.GetB1FromLocalStorage()
	xS := c.SignXsWithPlaywright(uri, data, method)
	xT := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return map[string]interface{}{
		"x-s":          xS,
		"x-t":          xT,
		"x-s-common":   BuildXsCommon(a1, b1, xS, xT),
		"x-b3-traceid": GetTraceId(),
	}
}

func (c *PlaywrightClient) PreHeadersWithPlaywright(
	urlStr string,
	cookieDict map[string]string,
	params map[string]string,
	payload map[string]interface{},
) map[string]string {
	a1Value := ""
	if a1, ok := cookieDict["a1"]; ok {
		a1Value = a1
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return map[string]string{}
	}
	uri := parsedURL.Path

	var data interface{}
	httpMethod := "POST"

	if params != nil && len(params) > 0 {
		data = params
		httpMethod = "GET"
	} else if payload != nil {
		data = payload
		httpMethod = "POST"
	}

	signs := c.SignWithPlaywright(uri, data, a1Value, httpMethod)

	return map[string]string{
		"X-S":          signs["x-s"].(string),
		"X-T":          signs["x-t"].(string),
		"x-S-Common":   signs["x-s-common"].(string),
		"X-B3-Traceid": signs["x-b3-traceid"].(string),
	}
}

func CustomQuote(s string, safeComma bool) string {
	result := strings.Builder{}
	for _, r := range []rune(s) {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if strings.ContainsRune("-_.~()*!.'", r) {
			result.WriteRune(r)
		} else if safeComma && r == ',' {
			result.WriteRune(r)
		} else if r == ' ' {
			result.WriteString("%20")
		} else {
			for _, b := range []byte(string(r)) {
				result.WriteString(fmt.Sprintf("%%%02X", b))
			}
		}
	}
	return result.String()
}
