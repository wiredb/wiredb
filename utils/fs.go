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
	"os"
)

// IsExist checked directory is exist
func IsExist(dirPath string) bool {
	// 使用 os.Stat 检查目录是否存在
	_, err := os.Stat(dirPath)

	// 如果 err 不为 nil 并且是目录不存在错误返回 false
	// 如果 err 为 nil 或者是其他类型的错误，权限问题则返回 true
	return !(err != nil && os.IsNotExist(err))
}

// IsDir check if the path is a directory
// !IsDir 就可以检查是否为文件
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// FlushToDisk 封装了文件的 Sync 和 Close 操作，减少重复代码
func FlushToDisk(fd *os.File) error {
	err := fd.Sync()
	if err != nil {
		return fmt.Errorf("failed to flush to disk: %w", err)
	}

	err = fd.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}

// BytesToGB converts a given size in bytes to gigabytes (GB).
func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
