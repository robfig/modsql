// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	. "github.com/kless/modsql"
	"time"
)

func main() {
	metadata := Metadata(PostgreSQL, MySQL, SQLite)

	types := Table("types", metadata,
		Column("t_int", Int).PrimaryKey(),
		Column("t_int8", Int8),
		Column("t_int16", Int16),
		Column("t_int32", Int32),
		Column("t_int64", Int64),

		Column("t_float32", Float32),
		Column("t_float64", Float64),

		Column("t_string", String).Unique(),
		Column("t_binary", Binary),

		Column("t_byte", Byte),
		Column("t_rune", Rune),
		Column("t_bool", Bool),
	)
	types.Unique("t_float32", "t_float64")

	def := Table("default_value", metadata,
		Column("id", Int).PrimaryKey(),
		Column("d_int8", Int8).Default(int8(55)),
		Column("d_float32", Float32).Default(float32(10.2)),

		Column("d_string", String),
		Column("d_binary", Binary),

		Column("d_byte", Byte).Default(byte('b')),
		Column("d_rune", Rune).Default('r'),
		Column("d_bool", Bool).Default(false),
	)

	times := Table("times", metadata,
		Column("typeId", Int).ForeignKey("types", "t_int"),
		Column("t_duration", Duration),
		Column("t_datetime", DateTime),
	)

	// Insert values

	types.Insert(
		1, 8, 16, 32, 64,
		1.32, 1.64,
		"one", []byte("12"),
		"A", "Z", true,
	)

	def.Insert(
		1, 10, 10.10,
		"foo", []byte{'1', '2'},
		"a", "z", false,
	)

	times.Insert(1, 5*time.Hour+3*time.Minute+12*time.Second, time.Now())

	// * * *

	metadata.Create().Write()
}
