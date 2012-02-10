// Copyright 2010  The "GoSQL" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gosql

// For columns with a wrong type
var (
	anyColumnErr bool
	columnsErr   []string
)

type column struct {
	name         string
	type_        sqlType
	isPrimaryKey bool
	defaultValue interface{}
}

func Column(name string, type_ sqlType) *column {
	col := new(column)
	col.name = name
	col.type_ = type_
	return col
}

func (c *column) Default(i interface{}) *column {
	c.defaultValue = i

	if ok := c.check(); !ok {
		columnsErr = append(columnsErr, c.name)
		anyColumnErr = true
	}

	return c
}

func (c *column) PrimaryKey() *column {
	c.isPrimaryKey = true
	return c
}

// Checks if the value by default has the correct type.
func (c *column) check() bool {
	switch c.defaultValue.(type) {
	case bool:
		if c.type_ != Boolean {
			return false
		}
	case float32, float64:
		if c.type_ != Float {
			return false
		}
	case int:
		if c.type_ != Integer {
			return false
		}
	case string:
		if c.type_ != Text {
			return false
		}
	default:
		return false
	}

	return true
}
