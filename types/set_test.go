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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestNewSet(t *testing.T) {
	set := NewSet()
	assert.NotNil(t, set)    // 确保 set 不为空
	assert.Empty(t, set.Set) // 确保新建的 set 为空
}

func TestSet_Add(t *testing.T) {
	set := NewSet()
	set.Add("apple")
	set.Add("banana")

	assert.True(t, set.Contains("apple"))
	assert.True(t, set.Contains("banana"))
	assert.False(t, set.Contains("orange")) // 不存在的元素
	assert.Equal(t, 2, set.Size())          // 确保 Set 的大小正确
}

func TestSet_Contains(t *testing.T) {
	set := NewSet()
	set.Add("grape")

	assert.True(t, set.Contains("grape"))
	assert.False(t, set.Contains("watermelon")) // 未添加的元素
}

func TestSet_Remove(t *testing.T) {
	set := NewSet()
	set.Add("apple")
	set.Add("banana")

	set.Remove("apple")                    // 删除元素
	assert.False(t, set.Contains("apple")) // 确保 apple 被删除
	assert.True(t, set.Contains("banana")) // 确保 banana 仍然存在
	assert.Equal(t, 1, set.Size())         // 确保 Set 的大小正确
}

func TestSet_Size(t *testing.T) {
	set := NewSet()
	assert.Equal(t, 0, set.Size()) // 空 Set

	set.Add("one")
	set.Add("two")
	assert.Equal(t, 2, set.Size()) // 添加两个元素

	set.Remove("one")
	assert.Equal(t, 1, set.Size()) // 删除一个元素
}

func TestSet_Clear(t *testing.T) {
	set := NewSet()
	set.Add("a")
	set.Add("b")
	set.TTL = 100

	set.Clear()                         // 清空
	assert.Equal(t, 0, set.Size())      // Set 应为空
	assert.Equal(t, uint64(0), set.TTL) // TTL 应重置为 0
}

func TestSet_ToBSON(t *testing.T) {
	set := NewSet()
	set.Add("x")
	set.Add("y")

	data, err := set.ToBSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, data) // 确保序列化后的 BSON 不为空

	// 反序列化回 Set 进行验证
	var decodedSet map[string]bool
	err = bson.Unmarshal(data, &decodedSet)
	assert.NoError(t, err)
	assert.Equal(t, set.Set, decodedSet) // 确保反序列化后的数据与原始数据一致
}
