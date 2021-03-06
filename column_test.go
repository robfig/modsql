// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"testing"
	"time"
)

func TestDefaultValue(t *testing.T) {
	val1 := false
	Column("married", Bool).Default(val1)
	checkError(t, val1)

	val2 := float32(12.2)
	Column("height", Float32).Default(val2)
	checkError(t, val2)

	val3 := int32(16)
	Column("age", Int32).Default(val3)
	checkError(t, val3)

	val4 := byte('a')
	Column("char", Byte).Default(val4)
	checkError(t, val4)

	//val5 := 10 * time.Second
	//Column("duration", Duration).Default(val5)
	//checkError(t, val5)

	val6 := time.Now()
	Column("time", DateTime).Default(val6)
	checkError(t, val6)
}

func TestConstraintWinthIndex(t *testing.T) {
	Column("foo", Int).PrimaryKey().Index(true)
	Column("bar", Bool).Index(false).Unique()

	if len(columnsErr) != 2 {
		t.Error("expected to get errors at set both constraint and index")
	}
}

// * * *

func checkError(t *testing.T, value interface{}) {
	if len(columnsErr) != 0 {
		t.Error("got error in column with value:", value)
	}
}
