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

func TestTextOperations(t *testing.T) {
	text := NewText("Hello, World!")
	assert.Equal(t, 13, text.Size())
	assert.True(t, text.Contains("World"))
	assert.False(t, text.Contains("Leon Ding"))

	text.Append(" and more")
	assert.Equal(t, "Hello, World! and more", text.Content)

	clone := text.Clone()
	assert.Equal(t, text.Content, clone.Content)
	assert.Equal(t, text.TTL, clone.TTL)

	text.Clear()
	assert.Equal(t, "", text.Content)
	assert.Equal(t, uint64(0), text.TTL)

	bsonData, err := text.ToBSON()
	assert.NoError(t, err)

	var decodedText Text
	err = bson.Unmarshal(bsonData, &decodedText)
	assert.NoError(t, err)
	assert.Equal(t, text, &decodedText)

	other := NewText("other text content.")
	assert.False(t, text.Equals(other))
}
