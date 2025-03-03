// Copyright 2022 Leon Ding <ding@ibyte.me> https://wiredb.github.io

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestSplitArgs(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"./cmd", "--port=2468", "--host=localhost", "--flag", "value"},
			expected: []string{"--port", "2468", "--host", "localhost", "--flag", "value"},
		},
		{
			input:    []string{"./cmd", "--port==8080", "--port===8080", "--flag=value"},
			expected: []string{"--flag", "value"},
		},
		{
			input:    []string{"./cmd", "arg1", "arg2", "arg3"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			result := SplitArgs(testCase.input)
			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("Expected %v, but got %v", testCase.expected, result)
			}
		})
	}
}

func TestTrimDaemon(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		// 测试移除 "-daemon" 参数的情况
		{
			input:    []string{"app", "-daemon", "arg1", "arg2", "--daemon", "arg3"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
		// 测试不包含 "-daemon" 参数的情况
		{
			input:    []string{"app", "arg1", "arg2", "arg3"},
			expected: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := TrimDaemon(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	length := 16
	var randomStr string

	for i := 0; i < length; i++ {
		randomStr = RandomString(length)
		t.Logf("Random String: %s", randomStr)
	}

	randomStr = RandomString(length + 1)

	// 检查随机字符串的长度是否正确
	if utf8.RuneCountInString(randomStr) != length {
		t.Errorf("Expected length %d, but got %d", length, utf8.RuneCountInString(randomStr))
	}
}
