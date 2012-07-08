// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package modsql

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
