// Copyright 2010  The "SQLModel" Authors
//
// Use of this source code is governed by the BSD-2 Clause license
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package sqlmodel

import (
	"testing"
)


func TestColumn(t *testing.T) {
	val1 := false
	Column("married", Boolean).Default(val1)
	checkError(t, val1)

	val2 := 12.2
	Column("height", Float).Default(val2)
	checkError(t, val2)

	val3 := 16
	Column("age", Integer).Default(val3)
	checkError(t, val3)

	val4 := "Pak"
	Column("name", Text).Default(val4)
	checkError(t, val4)

}

// ===

func checkError(t *testing.T, value interface{}) {
	if anyColumnErr == true {
		t.Error("It must have an error for:", value)
	}
	anyColumnErr = false
}

