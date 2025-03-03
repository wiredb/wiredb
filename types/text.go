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
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Text struct {
	Content string `json:"content" bson:"content" binding:"required"`
	TTL     uint64 `json:"ttl,omitempty"`
}

func NewText(content string) *Text {
	return &Text{Content: content}
}

func (text *Text) Size() int {
	return len(text.Content)
}

func (text *Text) Append(content string) {
	text.Content += content
}

func (text *Text) Contains(target string) bool {
	return strings.Contains(text.Content, target)
}

func (text *Text) Equals(other *Text) bool {
	return text.Content == other.Content
}

func (text *Text) Clone() *Text {
	return &Text{Content: text.Content, TTL: text.TTL}
}

func (text *Text) Clear() {
	text.TTL = 0
	text.Content = ""
}

func (text Text) ToBSON() ([]byte, error) {
	return bson.Marshal(text)
}
