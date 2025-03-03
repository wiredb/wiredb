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

func TestNumberOperations(t *testing.T) {
	num := NewNumber(10)
	assert.Equal(t, int64(10), num.Get())

	num.Add(5)
	assert.Equal(t, int64(15), num.Get())

	num.Sub(3)
	assert.Equal(t, int64(12), num.Get())

	num.Increment()
	assert.Equal(t, int64(13), num.Get())

	num.Decrement()
	assert.Equal(t, int64(12), num.Get())

	num.Set(100)
	assert.Equal(t, int64(100), num.Get())

	success := num.CompareAndSwap(100, 200)
	assert.True(t, success)
	assert.Equal(t, int64(200), num.Get())

	success = num.CompareAndSwap(100, 300)
	assert.False(t, success)
	assert.Equal(t, int64(200), num.Get())

	bsonData, err := num.ToBSON()
	assert.NoError(t, err)

	var decodedNumber Number
	err = bson.Unmarshal(bsonData, &decodedNumber)
	assert.NoError(t, err)
	assert.Equal(t, num.Get(), decodedNumber.Get())
}
