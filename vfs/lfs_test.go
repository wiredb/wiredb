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
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/auula/wiredb/conf"
	"github.com/auula/wiredb/types"
	"github.com/stretchr/testify/assert"
)

// TestSerializedIndex 测试 serializedIndex 函数
func TestSerializedIndex(t *testing.T) {
	// 创建一个测试的 INode
	inode := &INode{
		RegionID:  1234,
		Position:  5678,
		Length:    100,
		ExpiredAt: 1617181723,
		CreatedAt: 1617181623,
	}

	// 计算预期的字节切片
	expectedLength := 48

	// 调用 serializeIndex
	result, err := serializedIndex(1001, inode)
	if err != nil {
		t.Fatalf("serialized index failed: %v", err)
	}

	// 检查返回的字节切片长度
	if len(result) != expectedLength {
		t.Errorf("expected result length %d, got %d", expectedLength, len(result))
	}

	// 验证内容字段进行反序列化并检查
	inum, node, err := deserializedIndex(result)

	if err != nil {
		t.Errorf("failed to deserialized: %v", err)
	}

	// 验证字段是否一致
	if inum != 1001 {
		t.Errorf("expected inum %d, got %d", 1001, inum)
	}
	if node.RegionID != inode.RegionID {
		t.Errorf("expected RegionID %d, got %d", inode.RegionID, node.RegionID)
	}
	if node.Position != inode.Position {
		t.Errorf("expected Offset %d, got %d", inode.Position, node.RegionID)
	}
	if node.Length != inode.Length {
		t.Errorf("expected Length %d, got %d", inode.Length, node.Length)
	}
	if node.ExpiredAt != inode.ExpiredAt {
		t.Errorf("expected ExpiredAt %d, got %d", inode.ExpiredAt, node.ExpiredAt)
	}
	if node.CreatedAt != inode.CreatedAt {
		t.Errorf("expected CreatedAt %d, got %d", inode.CreatedAt, node.CreatedAt)
	}

}

// 测试 readSegment 函数
func TestReadSegment(t *testing.T) {
	// 构造测试数据
	seg := &Segment{
		Tombstone: 0,
		Type:      1,
		ExpiredAt: 123456789,
		CreatedAt: 987654321,
		KeySize:   3,
		ValueSize: 5,
		Key:       []byte("key"),
		Value:     []byte("value"),
	}

	// 将 Segment 数据转化为字节数组
	bytes, err := serializedSegment(seg)
	if err != nil {
		t.Fatalf("failed to serialized segment:%v", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试数据
	_, err = tmpFile.Write(bytes)
	if err != nil {
		t.Fatalf("failed to write test data to temp file: %v", err)
	}

	// 使用 readSegment 读取并测试数据
	offset := uint64(0)
	inum, segment, err := readSegment(tmpFile, offset, 26)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// 校验 Segment 数据
	if segment.Tombstone != seg.Tombstone {
		t.Errorf("expected Tombstone to be %d, but got: %d", seg.Tombstone, segment.Tombstone)
	}
	if segment.Type != seg.Type {
		t.Errorf("expected Type to be %d, but got: %d", seg.Type, segment.Type)
	}
	if segment.ExpiredAt != seg.ExpiredAt {
		t.Errorf("expected ExpiredAt to be %d, but got: %d", seg.ExpiredAt, segment.ExpiredAt)
	}
	if segment.CreatedAt != seg.CreatedAt {
		t.Errorf("expected CreatedAt to be %d, but got: %d", seg.CreatedAt, segment.CreatedAt)
	}
	if segment.KeySize != seg.KeySize {
		t.Errorf("expected KeySize to be %d, but got: %d", seg.KeySize, segment.KeySize)
	}
	if segment.ValueSize != seg.ValueSize {
		t.Errorf("expected ValueSize to be %d, but got: %d", seg.ValueSize, segment.ValueSize)
	}
	if string(segment.Key) != string(seg.Key) {
		t.Errorf("expected Key to be %s, but got: %s", string(seg.Key), string(segment.Key))
	}
	if string(segment.Value) != string(seg.Value) {
		t.Errorf("expected Value to be %s, but got: %s", string(seg.Value), string(segment.Value))
	}

	// 校验返回的 inode number (InodeNum)
	if inum != InodeNum(string(seg.Key)) {
		t.Errorf("expected InodeNum to be '%s', but got: %d", seg.Key, inum)
	}
}

func TestVFSWrite(t *testing.T) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`

	var tables types.Table
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		t.Fatal(err)
	}

	seg, err := NewSegment("key-01", tables, tables.TTL)
	if err != nil {
		t.Fatal(err)
	}

	err = fss.PutSegment("key-01", seg)
	if err != nil {
		t.Fatal(err)
	}

	_, seg, err = fss.FetchSegment("key-01")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", seg)
}

func BenchmarkVFSWrite(b *testing.B) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		b.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`
	var tables types.Table
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		b.Fatal(err)
	}

	// 重置计时器，忽略 setup 代码的影响
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)

		seg, err := NewSegment(key, tables, tables.TTL)
		if err != nil {
			b.Fatal(err)
		}

		err = fss.PutSegment(key, seg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVFSReads(b *testing.B) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, _, err := fss.FetchSegment(key)
		if err != nil {
			b.Logf("%v", err)
		}
	}
}

func TestUpdateSegmentWithCAS_Concurrent(t *testing.T) {
	var wg sync.WaitGroup

	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`
	var tables types.Table
	err = json.Unmarshal([]byte(data), &tables)
	if err != nil {
		t.Fatal(err)
	}

	key := "key-01"
	seg, err := NewSegment(key, tables, tables.TTL)
	if err != nil {
		t.Fatal(err)
	}

	err = fss.PutSegment(key, seg)
	if err != nil {
		t.Fatal(err)
	}

	var failures int32
	var success int32

	concurrentUpdates := rand.Intn(1000)
	for i := 0; i < concurrentUpdates; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 1. 获取当前版本号
			version, _, err := fss.FetchSegment(key)
			if err != nil {
				t.Errorf("goroutine %d: failed to fetch segment: %v", id, err)
				return
			}

			// 2. 创建新的 `Segment`
			newseg, err := NewSegment(key, tables, tables.TTL)
			if err != nil {
				t.Errorf("goroutine %d: failed to create segment: %v", id, err)
				return
			}

			// 3. CAS 更新
			err = fss.UpdateSegmentWithCAS(key, version, newseg)
			if err != nil {
				atomic.AddInt32(&failures, 1)
				t.Logf("goroutine %d: CAS update failed (expected version: %d)", id, version)
			} else {
				atomic.AddInt32(&success, 1)
				t.Logf("goroutine %d: CAS update succeeded (version: %d)", id, version)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Total success: %d, Total failures: %d,Updates concurrent: %d", success, failures, concurrentUpdates)

	// 断言至少有一些失败的情况（正常情况下应该有很多失败）
	if failures == 0 && (failures+success) != int32(concurrentUpdates) {
		t.Error("Expected some CAS failures, but got none")
	}
}

func TestConcurrentPutAndFetchSegment(t *testing.T) {
	var wg sync.WaitGroup
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 定义并发数量
	concurrentWrites := 100 // 写操作并发数
	concurrentReads := 100  // 读操作并发数

	wg.Add(concurrentWrites)
	// 存储测试数据
	for i := 0; i < concurrentWrites; i++ {
		go func(id int) {
			defer wg.Done()
			// 创建 segment
			k := fmt.Sprintf("key-%d", id)

			number := types.NewNumber(int64(id))

			segment, err := NewSegment(k, number, 0)
			if err != nil {
				t.Errorf("failed to create segment for key %s: %v", k, err)
				return
			}

			// 存储 segment
			err = fss.PutSegment(k, segment)
			if err != nil {
				t.Errorf("failed to put segment for key %s: %v", k, err)
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < concurrentReads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			k := fmt.Sprintf("key-%d", id)

			// 获取 segment
			_, seg, err := fss.FetchSegment(k)
			if err != nil {
				t.Errorf("failed to fetch segment for key: %s \t %v", k, err)
				return
			}

			// 转换为 number 并验证
			number, err := seg.ToNumber()
			if err != nil {
				t.Errorf("failed to convert segment to number for key %s \t %v", k, err)
			}

			t.Logf("K:%s\tV:%v", k, number.Get())
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}

func TestVFSOpertions(t *testing.T) {
	fss, err := OpenFS(&Options{
		FSPerm:    conf.FSPerm,
		Path:      conf.Settings.Path,
		Threshold: conf.Settings.Region.Threshold,
	})
	assert.NoError(t, err)

	data := `
{
  "table": {
    "is_valid": false,
    "items": [
      {
        "id": 1,
        "name": "Item 1"
      },
      {
        "id": 2,
        "name": "Item 2"
      }
    ],
    "meta": {
      "version": "2.0",
      "author": "Leon Ding"
    }
  }
}
`

	var tables types.Table
	err = json.Unmarshal([]byte(data), &tables)
	assert.NoError(t, err)

	seg, err := NewSegment("key-01", tables, tables.TTL)
	assert.NoError(t, err)

	err = fss.PutSegment("key-01", seg)
	assert.NoError(t, err)

	assert.Equal(t, fss.KeysCount(), 1)

	err = fss.DeleteSegment("key-01")
	assert.NoError(t, err)

	_, _, err = fss.FetchSegment("key-01")
	assert.Equal(t, err.Error(), "inode index for 9171687345308829835 not found")

	err = fss.ExportSnapshotIndex()
	assert.NoError(t, err)

	err = fss.SetEncryptor(AESCryptor, []byte("1234567890123456"))
	assert.NoError(t, err)

	fss.SetCompressor(SnappyCompressor)

	transformer = NewTransformer()

	os.RemoveAll(conf.Settings.Path)
}
