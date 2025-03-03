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
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsDirExist(t *testing.T) {
	// 测试存在的目录
	existingDir := os.TempDir()
	exists := IsExist(existingDir)
	if !exists {
		t.Errorf("Expected directory %s to exist, but it does not.", existingDir)
	}

	// 测试不存在的目录
	nonExistingDir := "/aaa/bbb/cccc/directory"
	exists = IsExist(nonExistingDir)
	if exists {
		t.Errorf("Expected directory %s to not exist, but it does.", nonExistingDir)
	}

	// 测试无效路径
	invalidPath := "/invalid/path"
	exists = IsExist(invalidPath)
	if exists {
		t.Errorf("Expected directory %s to not exist, but it does.", invalidPath)
	}
}

func TestIsDir(t *testing.T) {
	t.Run("Existing Directory", func(t *testing.T) {
		path := "."
		if !IsDir(path) {
			t.Errorf("Expected %s to be a directory", path)
		}
	})

	t.Run("Existing File", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file.Name())

		if IsDir(file.Name()) {
			t.Errorf("Expected %s to be a file, not a directory", file.Name())
		}
	})

	t.Run("Non-existent Path", func(t *testing.T) {
		if IsDir("/non/existent/path") {
			t.Errorf("Expected non-existent path to return false")
		}
	})
}

func TestFlushToDisk(t *testing.T) {
	invalidFd := os.NewFile(99999, "invalid") // 99999 是无效的文件描述符
	err := FlushToDisk(invalidFd)
	if err == nil {
		t.Errorf("Expected error for invalid file descriptor, but got nil")
	} else {
		t.Logf("Received expected error: %v", err)
	}
}

func TestFlushToDisk_SyncError_CloseError(t *testing.T) {
	// 创建临时文件
	tmpfile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file：%v", err)
	}

	// 确保在测试结束时删除临时文件
	defer os.Remove(tmpfile.Name())

	// 向临时文件写入数据
	if _, err := tmpfile.Write([]byte("Hello, World!")); err != nil {
		t.Fatalf("failed to write temp file：%v", err)
	}

	// 调用 FlushToDisk 函数
	if err := FlushToDisk(tmpfile); err != nil {
		t.Error(err)
	}

	// 检查临时文件是否已关闭
	if _, err := tmpfile.Write([]byte("test")); err != nil {
		t.Log(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Log(err)
	}
}

func TestBytesToGB(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected float64
	}{
		{1073741824, 1.0}, // 1 GB
		{2147483648, 2.0}, // 2 GB
		{536870912, 0.5},  // 0.5 GB
		{0, 0.0},          // 0 GB
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d bytes", tt.bytes), func(t *testing.T) {
			got := BytesToGB(tt.bytes)
			if got != tt.expected {
				t.Errorf("BytesToGB(%d) = %f; want %f", tt.bytes, got, tt.expected)
			}
		})
	}
}
