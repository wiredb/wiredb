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

package vfs

import (
	"testing"
)

// 测试 Transformer 类的压缩、加密和解密功能
func TestTransformerWithComplexData(t *testing.T) {
	// 创建一个新的 Transformer
	transformer := NewTransformer()

	// 构造复杂数据结构，包括 uint 和字符串
	originalString := "example-data"

	// 启用压缩
	transformer.SetCompressor(&Snappy{})
	transformer.SetEncryptor(AESCryptor, []byte("1234567890123456"))

	// 对数据进行编码（压缩 + 加密）
	encodedData, err := transformer.Encode([]byte(originalString))
	if err != nil {
		t.Fatalf("failed to encode data: %v", err)
	}

	// 解码数据
	decodedData, err := transformer.Decode(encodedData)
	if err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if originalString != string(decodedData) {
		t.Fatalf("failed to decode data: got %s, want %s", decodedData, originalString)
	}
}

// 测试 SnappyCompressor 类的压缩、加密和解密功能
func TestSnappyCompressor(t *testing.T) {
	// 构造复杂数据结构，包括 uint 和字符串
	originalString := "example-data"

	// 对数据进行编码（压缩 + 加密）
	encodedData, err := SnappyCompressor.Compress([]byte(originalString))
	if err != nil {
		t.Fatalf("failed to encode data: %v", err)
	}

	// 解码数据
	decodedData, err := SnappyCompressor.Decompress(encodedData)
	if err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if originalString != string(decodedData) {
		t.Fatalf("failed to decode data: got %s, want %s", decodedData, originalString)
	}
}

func TestCryptor(t *testing.T) {
	aes := new(Cryptor)

	// Example plaintext
	plaintext := []byte("Hello, this is a test of AES encryption!")

	// Key (must be either 16, 24, or 32 bytes long for AES-128, AES-192, or AES-256)
	secret := []byte("1234567890123456")

	// encrypt plaintext
	encrypted, err := aes.Encrypt(secret, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// decrypt ciphertext
	decrypted, err := aes.Decrypt(secret, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("got: %s , need: %s", decrypted, plaintext)
	}
}
