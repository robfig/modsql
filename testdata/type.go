// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import . "github.com/kless/modsql"

func main() {
	metadata := Metadata("mysql", Help)

	types := Table("types", metadata,
		Column("id", Int64).PrimaryKey(),
		Column("t_int", Int),
		Column("t_float", Float64),
		Column("t_string", String),
		//Column("t_binary", Binary),
		Column("t_bool", Bool),
	)

	def := Table("default_value", metadata,
		Column("id", Int64).PrimaryKey(),
		Column("d_int", Int8).Default(int8(55)),
		Column("d_float", Float32).Default(float32(10.2)),
		Column("d_string", String).Default("string"),
		//Column("d_binary", Binary).Default([]byte("123")),
		Column("d_bool", Bool).Default(false),
	)

	// == Insert values
	types.InsertHelp("en", "integer", "float", "text", "boolean")
	types.Insert(1, 10, 1.1, "one", true)
	types.Insert(2, 20, 2.2, "two", false)

	def.InsertHelp("en", "integer", "float", "text", "boolean")
	def.Insert(1, 10, 10.1, "foo", true)
	// ==

	metadata.Create().Write()
}
