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

import "gopkg.in/mgo.v2/bson"

type Set struct {
	Set map[string]bool `json:"set" bson:"set" binding:"required"`
	TTL uint64          `json:"ttl,omitempty"`
}

// 新建一个 Set
func NewSet() *Set {
	return &Set{
		Set: make(map[string]bool),
	}
}

// 向 Set 中添加一个元素
func (s *Set) Add(value string) {
	s.Set[value] = true
}

// 检查元素是否在 Set 中
func (s *Set) Contains(value string) bool {
	return s.Set[value]
}

// 从 Set 中删除一个元素
func (s *Set) Remove(value string) {
	delete(s.Set, value)
}

// 获取 Set 中的元素数量
func (s *Set) Size() int {
	return len(s.Set)
}

// 清空 Set
func (s *Set) Clear() {
	s.TTL = 0
	s.Set = make(map[string]bool)
}

func (s Set) ToBSON() ([]byte, error) {
	return bson.Marshal(s.Set)
}
