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
)

func TestList_AddItem(t *testing.T) {
	list := NewList()

	// Test adding an item
	item := "test item"
	list.AddItem(item)

	// Assert the item is added to the list
	assert.Equal(t, 1, list.Size())
	assert.Contains(t, list.List, item)
}

func TestList_Remove(t *testing.T) {
	list := NewList()
	item := "test item"
	list.AddItem(item)

	// Test removing an existing item
	err := list.Remove(item)
	assert.NoError(t, err)
	assert.NotContains(t, list.List, item)

	// Test removing a non-existing item
	err = list.Remove("non-existing item")
	assert.Error(t, err)
}

func TestList_GetItem(t *testing.T) {
	list := NewList()
	item := "test item"
	list.AddItem(item)

	// Test getting an item by index
	gotItem, err := list.GetItem(0)
	assert.NoError(t, err)
	assert.Equal(t, item, gotItem)

	// Test out-of-bounds index
	_, err = list.GetItem(1)
	assert.Error(t, err)
}

func TestList_Range(t *testing.T) {
	list := NewList()
	list.AddItem("item 1")
	list.AddItem("item 2")
	list.AddItem("item 3")

	// Test range function
	rangeItems, err := list.Rnage(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []any{"item 1", "item 2"}, rangeItems)

	// Test out-of-bounds range
	rangeItems, err = list.Rnage(2, 5)
	assert.NoError(t, err)
	assert.Equal(t, []any{"item 3"}, rangeItems)
}

func TestList_LPush(t *testing.T) {
	list := NewList()
	list.AddItem("item 1")
	list.LPush("new item")

	// Test LPush functionality
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, "new item", list.List[0])
}

func TestList_RPush(t *testing.T) {
	list := NewList()
	list.AddItem("item 1")
	list.RPush("new item")

	// Test RPush functionality
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, "new item", list.List[1])
}

func TestList_Clear(t *testing.T) {
	list := NewList()
	list.AddItem("item 1")
	list.Clear()

	// Test clear functionality
	assert.Equal(t, 0, list.Size())
	assert.Equal(t, uint64(0), list.TTL)
}

func TestList_ToBSON(t *testing.T) {
	list := NewList()
	list.AddItem("item 1")

	// Test ToBSON
	_, err := list.ToBSON()
	assert.NoError(t, err)
}
