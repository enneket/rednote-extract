package tools

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ParseResult struct {
	Scheme   string
	Netloc   string
	Path     string
	Params   string
	Query    string
	Fragment string
	Username string
	Password string
	Hostname string
	Port     string
}

type SplitResult struct {
	Scheme   string
	Netloc   string
	Path     string
	Query    string
	Fragment string
	Username string
	Password string
	Hostname string
	Port     string
}

func urlParse(rawUrl string) *ParseResult {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil
	}

	result := &ParseResult{
		Scheme:   u.Scheme,
		Netloc:   u.Host,
		Path:     u.Path,
		Query:    u.RawQuery,
		Fragment: u.Fragment,
	}

	if u.User != nil {
		result.Username = u.User.Username()
		result.Password, _ = u.User.Password()
	}

	result.Hostname = u.Hostname()
	result.Port = u.Port()

	return result
}

func urlSplit(rawUrl string) *SplitResult {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil
	}

	result := &SplitResult{
		Scheme:   u.Scheme,
		Netloc:   u.Host,
		Path:     u.Path,
		Query:    u.RawQuery,
		Fragment: u.Fragment,
	}

	if u.User != nil {
		result.Username = u.User.Username()
		result.Password, _ = u.User.Password()
	}

	result.Hostname = u.Hostname()
	result.Port = u.Port()

	return result
}

func urlUnparse(result ParseResult) string {
	u := &url.URL{
		Scheme:   result.Scheme,
		Path:     result.Path,
		RawQuery: result.Query,
		Fragment: result.Fragment,
	}

	if result.Netloc != "" || result.Username != "" {
		host := result.Netloc
		if result.Hostname != "" {
			host = result.Hostname
			if result.Port != "" {
				host = host + ":" + result.Port
			}
		}

		if result.Username != "" {
			userInfo := url.UserPassword(result.Username, result.Password)
			u.User = userInfo
			if result.Netloc != "" {
				u.Host = host
			}
		} else {
			u.Host = host
		}
	}

	return u.String()
}

func urlUnsplit(result SplitResult) string {
	u := &url.URL{
		Scheme:   result.Scheme,
		Path:     result.Path,
		RawQuery: result.Query,
		Fragment: result.Fragment,
	}

	if result.Netloc != "" || result.Username != "" {
		host := result.Netloc
		if result.Hostname != "" {
			host = result.Hostname
			if result.Port != "" {
				host = host + ":" + result.Port
			}
		}

		if result.Username != "" {
			userInfo := url.UserPassword(result.Username, result.Password)
			u.User = userInfo
			if result.Netloc != "" {
				u.Host = host
			}
		} else {
			u.Host = host
		}
	}

	return u.String()
}

func urlJoin(base, urlStr string) string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return urlStr
	}

	return baseUrl.ResolveReference(&url.URL{Path: urlStr}).String()
}

func parseQs(query string) map[string][]string {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return make(map[string][]string)
	}
	return parsed
}

type QPM struct {
	Key   string
	Value string
}

func parseQsl(query string) []QPM {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return []QPM{}
	}

	result := make([]QPM, 0, len(parsed))
	for key, values := range parsed {
		for _, value := range values {
			result = append(result, QPM{Key: key, Value: value})
		}
	}

	return result
}

func urlEncode(query map[string][]string) string {
	v := url.Values{}
	for key, values := range query {
		for _, value := range values {
			v.Add(key, value)
		}
	}
	return v.Encode()
}

func parseUrlEncoded(query string) map[string]string {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for key, values := range parsed {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

var hexDigits = "0123456789ABCDEFabcdef"

func Quote(s string, safe ...string) string {
	safeSet := make(map[byte]bool)
	for _, c := range ":/@!$&'()*+,;=" {
		safeSet[byte(c)] = true
	}
	for _, c := range safe {
		for i := 0; i < len(c); i++ {
			safeSet[c[i]] = true
		}
	}

	var builder strings.Builder
	builder.Grow(len(s) * 3)

	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
			builder.WriteByte(c)
		} else if safeSet[c] {
			builder.WriteByte(c)
		} else {
			builder.WriteByte('%')
			builder.WriteByte(hexDigits[c>>4])
			builder.WriteByte(hexDigits[c&0x0f])
		}
	}

	return builder.String()
}

func quoteForURI(s string) string {
	return Quote(s)
}

func QuotePlus(s string) string {
	s = strings.ReplaceAll(s, " ", "+")
	return Quote(s)
}

func unquote(s string) string {
	s, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return s
}

func unquotePlus(s string) string {
	s = strings.ReplaceAll(s, "+", " ")
	return unquote(s)
}

func unquoteToString(s string) string {
	return unquote(s)
}

func ExtractUrlParamsToDict(urlStr string) map[string]string {
	if urlStr == "" {
		return make(map[string]string)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return make(map[string]string)
	}

	query := u.RawQuery
	if query == "" {
		return make(map[string]string)
	}

	return parseUrlEncoded(query)
}

func urlencondeParams(params map[string]string) string {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value)
	}
	return v.Encode()
}

func base64UrlEncode(data []byte) string {
	result := base64.URLEncoding.EncodeToString(data)
	result = strings.ReplaceAll(result, "+", "-")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.TrimRight(result, "=")
	return result
}

func base64UrlDecode(s string) ([]byte, error) {
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")

	padding := len(s) % 4
	if padding > 0 {
		s += strings.Repeat("=", 4-padding)
	}

	return base64.URLEncoding.DecodeString(s)
}

func decodeBytes(s string) ([]byte, error) {
	return base64UrlDecode(s)
}

func decodeString(s string) (string, error) {
	data, err := base64UrlDecode(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func intToHex(i int) string {
	return strconv.FormatInt(int64(i), 16)
}

func hexToInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func urlParseParts(urlStr string) (string, string, string, string, string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", "", "", "", err
	}

	return u.Scheme, u.Host, u.Path, u.RawQuery, u.Fragment, nil
}

func urlunparseSimple(scheme, netloc, path, query, fragment string) string {
	u := &url.URL{
		Scheme:   scheme,
		Host:     netloc,
		Path:     path,
		RawQuery: query,
		Fragment: fragment,
	}
	return u.String()
}

func parseHTTPList(s string) []string {
	return strings.Fields(s)
}

func allowHosts(hosts []string, host string) bool {
	for _, h := range hosts {
		if h == host {
			return true
		}
	}
	return false
}

type Url struct {
	url.URL
}

func (u *Url) GetPort() string {
	return u.Port()
}

func (u *Url) GetHostname() string {
	return u.Hostname()
}

func (u *Url) GetAuthority() string {
	return u.Host
}

func (u *Url) GetUserInfo() string {
	if u.User == nil {
		return ""
	}
	return u.User.String()
}

func (u *Url) SetUserInfo(username, password string) {
	if username == "" {
		u.User = nil
	} else {
		u.User = url.UserPassword(username, password)
	}
}

func unquoteToBytes(s string) ([]byte, error) {
	result, err := url.QueryUnescape(s)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

type Quoter struct {
	safe string
}

func NewQuoter(safe string) *Quoter {
	return &Quoter{safe: safe}
}

func (q *Quoter) Quote(s string) string {
	return Quote(s, q.safe)
}

func (q *Quoter) unquote(s string) string {
	return unquote(s)
}

type Requoter struct {
	safe string
}

func NewRequoter(safe string) *Requoter {
	return &Requoter{safe: safe}
}

func (r *Requoter) Quote(s string) string {
	return Quote(s, r.safe)
}

func (r *Requoter) unquote(s string) string {
	return unquote(s)
}

func wrapRange(urlStr string, ranges []struct{ Start, End int }) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	u.RawQuery = "range=" + FormatRanges(ranges)
	return u.String()
}

func FormatRanges(ranges []struct{ Start, End int }) string {
	parts := make([]string, 0, len(ranges))
	for _, r := range ranges {
		parts = append(parts, strconv.Itoa(r.Start)+"-"+strconv.Itoa(r.End))
	}
	return strings.Join(parts, ",")
}

func UrlParseTarget(path string) (string, string, string, error) {
	idx := strings.Index(path, "?")
	if idx == -1 {
		return path, "", "", nil
	}

	path = path[:idx]
	query := path[idx+1:]

	idx = strings.Index(path, "#")
	if idx != -1 {
		query = path[idx+1:]
		path = path[:idx]
	}

	return path, query, "", nil
}

func isAbsolute(urlStr string) bool {
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}

func urlSplitQuery(urlStr string) (string, string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr, ""
	}

	return u.Path, u.RawQuery
}

func pathUrl2url(path string) string {
	u := &url.URL{Path: path}
	return u.String()
}

func unquotePlusToString(s string) string {
	return unquotePlus(s)
}

func unquoteToASCII(s string) string {
	return unquote(s)
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func parseAcceptEncoding(acceptEncoding string) map[string]float64 {
	encodings := make(map[string]float64)
	parts := strings.Split(acceptEncoding, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		idx := strings.Index(part, ";")
		if idx == -1 {
			encodings[part] = 1.0
		} else {
			encoding := strings.TrimSpace(part[:idx])
			qStr := strings.TrimSpace(part[idx+1:])
			if strings.HasPrefix(qStr, "q=") {
				q, err := strconv.ParseFloat(qStr[2:], 64)
				if err == nil {
					encodings[encoding] = q
				} else {
					encodings[encoding] = 1.0
				}
			} else {
				encodings[encoding] = 1.0
			}
		}
	}

	return encodings
}

func isSchemeValid(scheme string) bool {
	if len(scheme) == 0 {
		return false
	}

	for i, c := range scheme {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (i > 0 && ((c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.'))) {
			return false
		}
	}
	return true
}

func clearUrlCache() {
	url.Parse("")
}

func cacheableEncode(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	return Quote(s)
}

func cacheableDecode(s string) string {
	return unquote(s)
}

func sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
