package tools

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"
)

// GetRandomUserAgent 获取随机User-Agent
func GetRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0",
	}

	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[index.Int64()]
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// Base36Encode 将*big.Int编码为base36字符串
func Base36Encode(num int64) string {
	// 将int64转换为big.Int进行处理
	return Base36EncodeBigInt(big.NewInt(num))
}

// Base36EncodeBigInt 将*big.Int编码为base36字符串
func Base36EncodeBigInt(num *big.Int) string {
	if num == nil {
		return "0"
	}

	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result []byte
	base := big.NewInt(36)

	// 使用绝对值进行计算
	n := new(big.Int).Set(num)
	if n.Sign() < 0 {
		n.Abs(n)
	}

	zero := big.NewInt(0)
	temp := new(big.Int)
	remainder := new(big.Int)

	for n.Cmp(zero) > 0 {
		// 计算余数: n % base
		temp.DivMod(n, base, remainder)
		// 将余数转换为字符并添加到结果
		result = append([]byte{chars[remainder.Int64()]}, result...)
		// 更新n为商
		n.Set(temp)
	}

	// 如果结果为空，返回"0"
	if len(result) == 0 {
		return "0"
	}

	return string(result)
}

// SleepRandom 随机休眠指定时间范围内的秒数
func SleepRandom(min, max int) {
	if min >= max {
		time.Sleep(time.Duration(min) * time.Second)
		return
	}

	seconds, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	time.Sleep(time.Duration(min+int(seconds.Int64())) * time.Second)
}

// ParseURL 解析URL
func ParseURL(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// ExtractQueryParam 从URL中提取查询参数
func ExtractQueryParam(rawURL, key string) string {
	parsedURL, err := ParseURL(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Query().Get(key)
}

// GetSearchID 生成搜索ID
func GetSearchID() string {
	// 使用big.Int处理大整数，避免溢出
	// 计算当前时间戳（毫秒）
	timestamp := big.NewInt(time.Now().UnixMilli())
	// 左移64位
	e := new(big.Int).Lsh(timestamp, 64)
	// 生成0到2147483646之间的随机整数
	t, _ := rand.Int(rand.Reader, big.NewInt(2147483647)) // 2147483647 = 2^31 -1
	// 将两者相加
	total := new(big.Int).Add(e, t)
	// 使用Base36EncodeBigInt直接处理big.Int
	return Base36EncodeBigInt(total)
}

// ConvertCookies 将浏览器cookies转换为字符串
func ConvertCookies(cookies []map[string]interface{}) (string, map[string]string) {
	var cookieStr string
	cookieMap := make(map[string]string)

	for _, cookie := range cookies {
		if name, ok := cookie["name"].(string); ok {
			if value, ok := cookie["value"].(string); ok {
				cookieStr += fmt.Sprintf("%s=%s; ", name, value)
				cookieMap[name] = value
			}
		}
	}

	cookieStr = strings.TrimSuffix(cookieStr, "; ")
	return cookieStr, cookieMap
}

// IsValidURL 检查URL是否有效
func IsValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

// GetCurrentTimestamp 获取当前时间戳（秒）
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTimestampMS 获取当前时间戳（毫秒）
func GetCurrentTimestampMS() int64 {
	return time.Now().UnixMilli()
}
