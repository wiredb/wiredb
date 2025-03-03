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
	"math/rand"
	"strings"
)

// Charset defines the set of characters to be used in generating random strings
const Charset = "#$@!abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

// TrimDaemon removes the "-daemon" or "--daemon" arguments from os.Args
func TrimDaemon(args []string) []string {
	var newArgs []string

	// Iterate through the args slice
	for i := 1; i < len(args); i++ {
		// Skip the current argument if it matches the daemon flags
		if args[i] == "-daemon" || args[i] == "--daemon" {
			continue
		}
		newArgs = append(newArgs, args[i])
	}

	return newArgs
}

// SplitArgs splits command-line arguments by "=" if present
func SplitArgs(args []string) []string {
	var newArgs []string

	for i := 1; i < len(args); i++ {
		// Split elements in args by "=" to ensure proper command-line parsing
		if strings.Contains(args[i], "=") && strings.Count(args[i], "=") == 1 {
			newArgs = append(newArgs, strings.Split(args[i], "=")...)
		} else {
			// Skip elements with multiple "=" as they are invalid
			if strings.Count(args[i], "=") > 1 {
				continue
			}
			// ./cmd agrs1 agrs2 --debug --daemon
			newArgs = append(newArgs, args[i])
		}
	}

	return newArgs
}

// RandomString returns a string of the specified length composed of characters from Charset
func RandomString(length int) string {
	result := make([]byte, length-1)
	for i := 0; i < length-1; i++ {
		result[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(result)
}
