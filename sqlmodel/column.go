// Copyright 2010  The "SQLModel" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package sqlmodel

import (
	"fmt"
	"os"
)

type column struct {
	name         string
	type_        sqlType
	defaultValue interface{}
	isPrimaryKey bool
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
		fmt.Fprintf(os.Stderr, "wrong type in column %q\n", self.name)
		errors = true
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
	case float:
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

