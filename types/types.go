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

type QueryRow struct {
	Key   string `json:"key"`
	Index int    `json:"index,omitempty"`
	Range string `json:"range,omitempty"`
	Order string `json:"order,omitempty"`
	Score string `json:"score,omitempty"`
	Field string `json:"field,omitempty"`
}

type Query struct {
	Type  string  `json:"type"`
	Query []Query `json:"query"`
}

type Queryer interface {
	Search(qs []Query) []byte
}

type StrQuery struct {
	Query
}

func (sq *StrQuery) Search(qs []Query) []byte {
	return []byte{}
}

type HashQuery struct {
	Query
}

type ListQuery struct {
	Query
}

type SetQuery struct {
	Query
}

type ZSetQuery struct {
	Query
}
