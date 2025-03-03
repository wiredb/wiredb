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
	"errors"
	"math"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

type Double struct {
	Value float64 `json:"double" bson:"double" binding:"required"`
	TTL   uint64  `json:"ttl,omitempty"`
}

func NewDouble(num float64) *Double {
	return &Double{Value: num}
}

// 加法运算
func (d *Double) Add(x float64) {
	d.Value += x
}

// 减法运算
func (d *Double) Subtract(x float64) {
	d.Value -= x
}

// 乘法运算
func (d *Double) Multiply(x float64) {
	d.Value *= x
}

// 除法运算（防止除零错误）
func (d *Double) Divide(x float64) error {
	if x == 0 {
		return errors.New("division by zero")
	}
	d.Value /= x
	return nil
}

// 比较是否相等（考虑浮点数精度问题）
func (d *Double) Equals(x float64) bool {
	return math.Abs(d.Value-x) < 1e-9
}

// 判断是否大于
func (d *Double) GreaterThan(x float64) bool {
	return d.Value > x
}

// 判断是否小于
func (d *Double) LessThan(x float64) bool {
	return d.Value < x
}

// 舍入到指定小数位
func (d *Double) Round(precision int) {
	factor := math.Pow(10, float64(precision))
	d.Value = math.Round(d.Value*factor) / factor
}

// 转换为字符串
func (d *Double) String() string {
	return strconv.FormatFloat(d.Value, 'f', -1, 64)
}

// BSON 序列化
func (d Double) ToBSON() ([]byte, error) {
	return bson.Marshal(d)
}
