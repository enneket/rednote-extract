package tools

import (
	"strings"
	"testing"
)

func TestMrc(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "长度为57的普通字符串",
			input:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ12345",
			expected: -2926172971,
		},
		{
			name:     "包含特殊字符的字符串",
			input:    "!@#$%^&*()_+-=[]{}|;:,.<>?/~`'\"abcdefghijklmnopqrstuvwxyz",
			expected: -3167215328,
		},
		{
			name:     "全数字字符串",
			input:    "123456789012345678901234567890123456789012345678901234567",
			expected: -87372336,
		},
		{
			name:     "全大写字母字符串",
			input:    "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDE",
			expected: -4178307715,
		},
		{
			name:     "全小写字母字符串",
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcde",
			expected: -1249730137,
		},
		{
			name:     "包含重复字符的字符串",
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: -1110711967,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mrc(tt.input)
			if result != tt.expected {
				t.Errorf("Mrc(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTripletToBase64(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    uint32
		expected string
	}{
		{
			name:     "全零输入",
			input:    0,
			expected: "AAAA",
		},
		{
			name:     "小整数输入",
			input:    1,
			expected: "AAAB",
		},
		{
			name:     "大整数输入",
			input:    16777215,
			expected: "////",
		},
		{
			name:     "包含高字节的输入",
			input:    16711680,
			expected: "/wAA",
		},
		{
			name:     "包含中字节的输入",
			input:    65280,
			expected: "AP8A",
		},
		{
			name:     "包含低字节的输入",
			input:    255,
			expected: "AAD/",
		},
		{
			name:     "随机整数输入",
			input:    1193046,
			expected: "EjRW",
		},
		{
			name:     "ASCII字符组合 ('ABC')",
			input:    4276803,
			expected: "QUJD",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tripletToBase64(tt.input)
			if result != tt.expected {
				t.Errorf("tripletToBase64(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeChunk(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    []int
		start    int
		end      int
		expected string
	}{
		{
			name:     "3字节输入 (正好一个块)",
			input:    []int{0x41, 0x42, 0x43}, // 'ABC'
			start:    0,
			end:      3,
			expected: "QUJD",
		},
		{
			name:     "6字节输入 (两个块)",
			input:    []int{0x41, 0x42, 0x43, 0x44, 0x45, 0x46}, // 'ABCDEF'
			start:    0,
			end:      6,
			expected: "QUJDREVG",
		},
		{
			name:     "9字节输入 (三个块)",
			input:    []int{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49}, // 'ABCDEFGHI'
			start:    0,
			end:      9,
			expected: "QUJDREVGR0hJ",
		},
		{
			name:     "从中间开始的输入",
			input:    []int{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49}, // 'ABCDEFGHI'
			start:    3,
			end:      9,
			expected: "REVGR0hJ",
		},
		{
			name:     "全零字节输入",
			input:    []int{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // '000000000000'
			start:    0,
			end:      6,
			expected: "AAAAAAAA",
		},
		{
			name:     "包含特殊值的字节输入",
			input:    []int{0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF}, // 'ff000000ff000000ff'
			start:    0,
			end:      9,
			expected: "/wAAAP8AAAD/",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeChunk(tt.input, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("encodeChunk(%v, %d, %d) = %q, want %q", tt.input, tt.start, tt.end, result, tt.expected)
			}
		})
	}
}

func TestRightShiftUnsigned(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    int32
		shift    uint
		expected uint32
	}{
		{
			name:     "正数，位移0位",
			input:    12345,
			shift:    0,
			expected: 12345,
		},
		{
			name:     "正数，位移8位",
			input:    12345,
			shift:    8,
			expected: 48,
		},
		{
			name:     "正数，位移16位",
			input:    12345,
			shift:    16,
			expected: 0,
		},
		{
			name:     "正数，位移31位",
			input:    12345,
			shift:    31,
			expected: 0,
		},
		{
			name:     "负数，位移0位",
			input:    -1,
			shift:    0,
			expected: 4294967295,
		},
		{
			name:     "负数，位移8位",
			input:    -1,
			shift:    8,
			expected: 16777215,
		},
		{
			name:     "负数，位移16位",
			input:    -1,
			shift:    16,
			expected: 65535,
		},
		{
			name:     "负数，位移31位",
			input:    -1,
			shift:    31,
			expected: 1,
		},
		{
			name:     "其他负数，位移4位",
			input:    -12345,
			shift:    4,
			expected: 268434684,
		},
		{
			name:     "边界值0，位移10位",
			input:    0,
			shift:    10,
			expected: 0,
		},
		{
			name:     "最大32位无符号整数，位移16位",
			input:    int32(-1), // 等价于uint32的4294967295
			shift:    16,
			expected: 65535,
		},
		{
			name:     "最小32位有符号整数，位移8位",
			input:    -2147483648,
			shift:    8,
			expected: 8388608,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 将int32转换为uint32，模拟无符号处理
			result := RightShiftUnsigned(uint32(tt.input), tt.shift)
			if result != tt.expected {
				t.Errorf("RightShiftUnsigned(%d, %d) = %d, want %d", tt.input, tt.shift, result, tt.expected)
			}
		})
	}
}

func TestEncodeUtf8(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "普通ASCII字符",
			input:    "hello",
			expected: []int{104, 101, 108, 108, 111}, // [0x68, 0x65, 0x6C, 0x6C, 0x6F]
		},
		{
			name:     "包含空格的字符串",
			input:    "hello world",
			expected: []int{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}, // [0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64]
		},
		{
			name:     "包含特殊字符的字符串",
			input:    "hello@world#123",
			expected: []int{104, 101, 108, 108, 111, 64, 119, 111, 114, 108, 100, 35, 49, 50, 51}, // [0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x40, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x23, 0x31, 0x32, 0x33]
		},
		{
			name:     "包含中文字符的字符串",
			input:    "你好世界",
			expected: []int{228, 189, 160, 229, 165, 189, 228, 184, 150, 231, 149, 140}, // [0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C]
		},
		{
			name:     "包含混合字符的字符串",
			input:    "Hello 世界!@#123",
			expected: []int{72, 101, 108, 108, 111, 32, 228, 184, 150, 231, 149, 140, 33, 64, 35, 49, 50, 51}, // [0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0x21, 0x40, 0x23, 0x31, 0x32, 0x33]
		},
		{
			name:     "包含URL需要编码的字符",
			input:    "hello world&test=123",
			expected: []int{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 38, 116, 101, 115, 116, 61, 49, 50, 51}, // [0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x26, 0x74, 0x65, 0x73, 0x74, 0x3D, 0x31, 0x32, 0x33]
		},
		{
			name:     "包含安全字符的字符串",
			input:    "hello~()*!.'world",
			expected: []int{104, 101, 108, 108, 111, 126, 40, 41, 42, 33, 46, 39, 119, 111, 114, 108, 100}, // [0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x7E, 0x28, 0x29, 0x2A, 0x21, 0x2E, 0x27, 0x77, 0x6F, 0x72, 0x6C, 0x64]
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []int{}, // 空列表
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeUtf8(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("EncodeUtf8(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("EncodeUtf8(%q)[%d] = %d, want %d", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}
func TestCustomQuote(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "普通ASCII字符",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "包含空格的字符串",
			input:    "hello world",
			expected: "hello+world", // 注意：customQuote将空格替换为'+'
		},
		{
			name:     "包含特殊字符的字符串",
			input:    "hello@world#123",
			expected: "hello%40world%23123",
		},
		{
			name:     "包含中文字符的字符串",
			input:    "你好世界",
			expected: "%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C",
		},
		{
			name:     "包含URL保留字符的字符串",
			input:    "hello&world=123",
			expected: "hello%26world%3D123",
		},
		{
			name:     "包含安全字符的字符串",
			input:    "hello~()*!.",
			expected: "hello~()*!.",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customQuote(tt.input)
			if result != tt.expected {
				t.Errorf("customQuote(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestB64Encode(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    []int
		expected string
	}{
		{
			name:     "空列表",
			input:    []int{},
			expected: "",
		},
		{
			name:     "3字节输入 (正好一个块，无余数)",
			input:    []int{0x41, 0x42, 0x43}, // 'ABC'
			expected: "QUJD",
		},
		{
			name:     "1字节输入 (余数为1)",
			input:    []int{0x41}, // 'A'
			expected: "QQ==",
		},
		{
			name:     "2字节输入 (余数为2)",
			input:    []int{0x41, 0x42}, // 'AB'
			expected: "QUI=",
		},
		{
			name:     "4字节输入 (1个完整块 + 余数1)",
			input:    []int{0x41, 0x42, 0x43, 0x44}, // 'ABCD'
			expected: "QUJDRA==",
		},
		{
			name:     "5字节输入 (1个完整块 + 余数2)",
			input:    []int{0x41, 0x42, 0x43, 0x44, 0x45}, // 'ABCDE'
			expected: "QUJDREU=",
		},
		{
			name:     "包含中文字符的输入 (UTF-8编码)",
			input:    []int{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD}, // '你好'
			expected: "5L2g5aW9",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := B64Encode(tt.input)
			if result != tt.expected {
				t.Errorf("B64Encode(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}

	// 测试超过16383字节的情况
	t.Run("超过16383字节的情况", func(t *testing.T) {
		// 创建16384字节的输入数据
		largeInput := make([]int, 16384)
		for i := 0; i < 16384; i++ {
			largeInput[i] = i % 256 // 填充0-255的循环值
		}
		result := B64Encode(largeInput)
		// 验证输出长度是否符合预期 (16384 * 4 / 3 = 21845.333..., 向上取整并考虑填充)
		expectedLength := 21848
		if len(result) != expectedLength {
			t.Errorf("B64Encode large input length = %d, want %d", len(result), expectedLength)
		}
	})
}

func TestGetTraceId(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name string
	}{
		{
			name: "基本功能测试",
		},
		{
			name: "多次调用唯一性测试",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTraceId()

			// 验证返回值长度为16
			if len(result) != 16 {
				t.Errorf("GetTraceId() 长度 = %d, want %d", len(result), 16)
			}

			// 验证返回值只包含允许的字符 (abcdef0123456789)
			allowedChars := "abcdef0123456789"
			for _, char := range result {
				if !strings.ContainsRune(allowedChars, char) {
					t.Errorf("GetTraceId() 包含不允许的字符 %c", char)
				}
			}

			// 对于"多次调用唯一性测试"，验证多次生成的ID不相同
			if tt.name == "多次调用唯一性测试" {
				ids := make(map[string]bool)
				for i := 0; i < 100; i++ { // 生成100个ID来测试唯一性
					id := GetTraceId()
					if ids[id] {
						t.Errorf("GetTraceId() 生成了重复的ID: %s", id)
					}
					ids[id] = true

					// 同时验证每个ID的格式
					if len(id) != 16 {
						t.Errorf("GetTraceId() 长度 = %d, want %d", len(id), 16)
					}
					for _, char := range id {
						if !strings.ContainsRune(allowedChars, char) {
							t.Errorf("GetTraceId() 包含不允许的字符 %c", char)
						}
					}
				}
			}
		})
	}
}
