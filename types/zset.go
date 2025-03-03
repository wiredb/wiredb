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
	"sort"

	"gopkg.in/mgo.v2/bson"
)

// ZSet 是一个实现有序集合的结构
type ZSet struct {
	ZSet         map[string]float64 `json:"zset" bson:"zset" binding:"required"`
	TTL          uint64             `json:"ttl,omitempty"`
	sortedScores []string
}

// NewZSet 创建一个新的 ZSet
func NewZSet() *ZSet {
	return &ZSet{
		ZSet:         make(map[string]float64),
		sortedScores: []string{},
	}
}

// Add 向 ZSet 中添加一个元素，并指定它的分数
func (z *ZSet) Add(value string, score float64) {
	// 如果元素已经存在，更新其分数
	if _, exists := z.ZSet[value]; exists {
		z.Remove(value) // 删除旧的元素
	}
	// 将元素添加到 map 中
	z.ZSet[value] = score
	// 将元素添加到 sortedScores 列表中，并按分数排序
	z.sortedScores = append(z.sortedScores, value)
	z.sort()
}

// Remove 从 ZSet 中删除一个元素
func (z *ZSet) Remove(value string) {
	if _, exists := z.ZSet[value]; exists {
		delete(z.ZSet, value)
		// 更新 sortedScores 列表
		for i, v := range z.sortedScores {
			if v == value {
				z.sortedScores = append(z.sortedScores[:i], z.sortedScores[i+1:]...)
				break
			}
		}
	}
	z.sort()
}

// Get 获取元素的分数
func (z *ZSet) Get(value string) (float64, bool) {
	score, exists := z.ZSet[value]
	return score, exists
}

// GetRank 获取元素的排名（按分数排序）
func (z *ZSet) GetRank(value string) (int, bool) {
	z.sort()
	for i, v := range z.sortedScores {
		if v == value {
			return i, true
		}
	}
	return -1, false
}

// GetRange 获取指定分数区间内的元素
func (z *ZSet) GetRange(minScore, maxScore float64) []string {
	z.sort()
	var result []string
	for _, value := range z.sortedScores {
		if score, exists := z.ZSet[value]; exists && score >= minScore && score <= maxScore {
			result = append(result, value)
		}
	}
	return result
}

// sort 根据分数对 sortedScores 排序
func (z *ZSet) sort() {
	sort.Slice(z.sortedScores, func(i, j int) bool {
		return z.ZSet[z.sortedScores[i]] > z.ZSet[z.sortedScores[j]]
	})
}

func (z *ZSet) Size() int {
	return len(z.ZSet)
}

func (z *ZSet) Clear() {
	z.TTL = 0
	z.ZSet = make(map[string]float64)
	z.sortedScores = make([]string, 0)
}

func (zs ZSet) ToBSON() ([]byte, error) {
	return bson.Marshal(zs.ZSet)
}
