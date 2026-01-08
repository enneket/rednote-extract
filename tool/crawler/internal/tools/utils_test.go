package tools

import "testing"

func TestBase36Encode(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "输入0",
			input:    0,
			expected: "0",
		},
		{
			name:     "输入10",
			input:    10,
			expected: "A",
		},
		{
			name:     "输入36",
			input:    36,
			expected: "10",
		},
		{
			name:     "输入-10",
			input:    -10,
			expected: "-A",
		},
		{
			name:     "输入123456",
			input:    123456,
			expected: "2N9C",
		},
		{
			name:     "输入1000000000",
			input:    1000000000,
			expected: "GJDGXS",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base36Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Base36Encode(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
