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
	"time"

	"github.com/auula/wiredb/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	// Test valid Set type
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}

	// Create a new segment for the Set type
	segment, err := NewSegment("mock-key", set, 1000)
	assert.NoError(t, err)                                    // Ensure no error
	assert.NotNil(t, segment)                                 // Ensure segment is created
	assert.Equal(t, "mock-key", string(segment.Key))          // Ensure the key is set correctly
	assert.Equal(t, uint32(len("mock-key")), segment.KeySize) // Ensure the key size is correct
	assert.Equal(t, uint32(21), segment.ValueSize)            // Ensure the value size is correct
}

func TestNewTombstoneSegment(t *testing.T) {
	// Create a Tombstone segment
	segment := NewTombstoneSegment("mock-key")

	// Ensure the segment is of Tombstone type and has expected fields
	assert.Equal(t, Unknown, segment.Type)                    // Tombstone should have Unknown type
	assert.Equal(t, int8(1), segment.Tombstone)               // Tombstone should be marked as 1
	assert.Equal(t, "mock-key", string(segment.Key))          // Ensure the key is set correctly
	assert.Equal(t, uint32(len("mock-key")), segment.KeySize) // Ensure the key size is correct
}

func TestSegmentSize(t *testing.T) {
	// Create a Set type data for testing
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}

	// Create a segment for the Set type
	segment, err := NewSegment("mock-key", set, 1000)
	assert.NoError(t, err)

	// Ensure the size is calculated correctly
	assert.Equal(t, uint32(59), segment.Size())
}

func TestToSet(t *testing.T) {
	// Create a Set type Segment
	setData := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
		TTL: uint64(0),
	}
	segment, err := NewSegment("mock-key", setData, 1000)
	assert.NoError(t, err)

	// Convert the segment to Set
	set, err := segment.ToSet()
	assert.NoError(t, err)                // Ensure no error
	assert.Equal(t, setData.Set, set.Set) // Ensure the Set values match
}

func TestTTL(t *testing.T) {
	// Create a Segment with TTL
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}
	segment, err := NewSegment("mock-key", set, 1) // TTL = 1 second
	assert.NoError(t, err)

	// Wait 1 second
	time.Sleep(time.Second)

	// Test TTL, it should return a value close to 0
	ttl := segment.TTL()
	assert.True(t, ttl <= 0) // Ensure TTL is <= 0 after expiration
}

// TestToZSet 测试 ToZSet 方法
func TestToZSet(t *testing.T) {
	// 创建 ZSet 数据
	zsetData := types.ZSet{
		ZSet: map[string]float64{
			"user1": 100.5,
			"user2": 200.0,
		},
	}

	segment, err := NewSegment("test-key-01", zsetData, 0)
	assert.NoError(t, err)

	// 测试 ToZSet 方法
	result, err := segment.ToZSet()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, zsetData.ZSet, result.ZSet)
}

// TestToText 测试 ToText 方法
func TestToText(t *testing.T) {
	// 创建 Text 数据
	textData := types.Text{
		Content: "Hello, World!",
	}

	segment, err := NewSegment("test-key-01", textData, 0)
	assert.NoError(t, err)

	// 测试 ToText 方法
	result, err := segment.ToText()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, textData.Content, result.Content)
}

// TestToList 测试 ToList 方法
func TestToList(t *testing.T) {
	// 创建 List 数据
	listData := types.List{
		List: []any{"item1", "item2", 123},
	}

	segment, err := NewSegment("test-key-01", listData, 0)
	assert.NoError(t, err)

	// 测试 ToList 方法
	result, err := segment.ToList()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, listData.List, result.List)
}

// TestToTable 测试 ToTable 方法
func TestToTable(t *testing.T) {
	// 创建 Tables 数据
	tablesData := types.Table{
		Table: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	segment, err := NewSegment("test-key-01", tablesData, 0)
	assert.NoError(t, err)

	// 测试 ToTable 方法
	result, err := segment.ToTable()
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, tablesData.Table, result.Table)
}
