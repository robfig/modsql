// Copyright 2010  The "GotoSQL" Authors
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

package tosql

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

func (self *column) Default(i interface{}) *column {
	self.defaultValue = i

	if ok := self.check(); !ok {
		columnsErr = append(columnsErr, self.name)
		anyColumnErr = true
	}

	return self
}

func (self *column) PrimaryKey() *column {
	self.isPrimaryKey = true
	return self
}

// Checks if the value by default has the correct type.
func (self *column) check() bool {
	switch t := self.defaultValue.(type) {
	case bool:
		if self.type_ != Boolean {
			return false
		}
	case float32, float64:
		if self.type_ != Float {
			return false
		}
	case int:
		if self.type_ != Integer {
			return false
		}
	case string:
		if self.type_ != Text {
			return false
		}
	default:
		return false
	}

	return true
}
