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

import "testing"

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

// * * *

func checkError(t *testing.T, value interface{}) {
	if anyColumnErr == true {
		t.Error("It must have an error for:", value)
	}
	anyColumnErr = false
}
