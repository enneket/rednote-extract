package tools

import (
	"testing"
)

func TestMd5Hex(t *testing.T) {
	// 测试用例：空字符串
	if result := md5Hex(""); result != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("md5Hex(\"\") = %q, want %q", result, "d41d8cd98f00b204e9800998ecf8427e")
	}

	// 测试用例：简单字符串
	if result := md5Hex("hello"); result != "5d41402abc4b2a76b9719d911017c592" {
		t.Errorf("md5Hex(\"hello\") = %q, want %q", result, "5d41402abc4b2a76b9719d911017c592")
	}

	// 测试用例：简单字符串
	if result := md5Hex("hello world"); result != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Errorf("md5Hex(\"hello world\") = %q, want %q", result, "5eb63bbbe01eeed093cb22bb8f5acdc3")
	}

	// 测试用例：中文字符串
	if result := md5Hex("你好世界"); result != "65396ee4aad0b4f17aacd1c6112ee364" {
		t.Errorf("md5Hex(\"你好世界\") = %q, want %q", result, "65396ee4aad0b4f17aacd1c6112ee364")
	}

	if result := md5Hex("hello@world#123"); result != "5fd556a7a174fb95993f5eb391f0bcad" {
		t.Errorf("md5Hex(\"hello@world#123\") = %q, want %q", result, "5fd556a7a174fb95993f5eb391f0bcad")
	}
}
func TestBuildXSPayload(t *testing.T) {
	// 测试用例1：x3_value='test_value', data_type='object'
	if result := buildXSPayload("test_value", "object"); result != "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6InRlc3RfdmFsdWUiLCJ4NCI6Im9iamVjdCJ9" {
		t.Errorf("buildXSPayload(\"test_value\", \"object\") = %q, want %q", result, "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6InRlc3RfdmFsdWUiLCJ4NCI6Im9iamVjdCJ9")
	}

	// 测试用例2：x3_value='中文测试', data_type='object'
	if result := buildXSPayload("中文测试", "object"); result != "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6Ilx1NGUyZFx1NjU4N1x1NmQ0Ylx1OGJkNSIsIng0Ijoib2JqZWN0In0=" {
		t.Errorf("buildXSPayload(\"中文测试\", \"object\") = %q, want %q", result, "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6Ilx1NGUyZFx1NjU4N1x1NmQ0Ylx1OGJkNSIsIng0Ijoib2JqZWN0In0=")
	}

	// 测试用例3：x3_value='special@chars#123', data_type='string'
	if result := buildXSPayload("special@chars#123", "string"); result != "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6InNwZWNpYWxAY2hhcnMjMTIzIiwieDQiOiJzdHJpbmcifQ==" {
		t.Errorf("buildXSPayload(\"special@chars#123\", \"string\") = %q, want %q", result, "XYS_eyJ4MCI6IjQuMi4xIiwieDEiOiJ4aHMtcGMtd2ViIiwieDIiOiJNYWMgT1MiLCJ4MyI6InNwZWNpYWxAY2hhcnMjMTIzIiwieDQiOiJzdHJpbmcifQ==")
	}
}

func TestBuildSignString(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		data     interface{}
		method   string
		expected string
	}{
		{
			name:     "POST无数据",
			uri:      "/api/v1/test",
			data:     nil,
			method:   "POST",
			expected: "/api/v1/test",
		},
		{
			name: "POST带字典数据",
			uri:  "/api/v1/user",
			data: func() *OrderedMap {
				om := NewOrderedMap()
				om.Set("name", "张三")
				om.Set("age", 30)
				om.Set("active", true)
				return om
			}(),
			method:   "POST",
			expected: "/api/v1/user{\"name\":\"张三\",\"age\":30,\"active\":true}",
		},
		{
			name:     "POST带字符串数据",
			uri:      "/api/v1/raw",
			data:     "rawdata123",
			method:   "POST",
			expected: "/api/v1/rawrawdata123",
		},
		{
			name:     "GET无数据",
			uri:      "/api/v1/list",
			data:     nil,
			method:   "GET",
			expected: "/api/v1/list",
		},
		{
			name:     "GET带空字典",
			uri:      "/api/v1/empty",
			data:     map[string]interface{}{},
			method:   "GET",
			expected: "/api/v1/empty",
		},
		{
			name: "GET带多种类型字典数据",
			uri:  "/api/v1/query",
			data: func() *OrderedMap {
				om := NewOrderedMap()
				om.Set("name", "test user")
				om.Set("ids", []interface{}{1, 2, 3})
				om.Set("status", nil)
				om.Set("score", 95.5)
				return om
			}(),
			method:   "GET",
			expected: "/api/v1/query?name=test%20user&ids=1%2C2%2C3&status=&score=95.5",
		},
		{
			name:     "GET带字符串数据",
			uri:      "/api/v1/search",
			data:     "q=test&page=1",
			method:   "GET",
			expected: "/api/v1/search?q=test&page=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSignString(tt.uri, tt.data, tt.method)
			if result != tt.expected {
				t.Errorf("buildSignString(%q, %v, %q) = %q, want %q", tt.uri, tt.data, tt.method, result, tt.expected)
			}
		})
	}
}
