// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import . "github.com/kless/modsql"

func main() {
	metadata := Metadata(Help, PostgreSQL, MySQL, SQLite)

	types := Table("types", metadata,
		Column("t_int", Int).PrimaryKey(),
		Column("t_int8", Int8),
		Column("t_int16", Int16),
		Column("t_int32", Int32),
		Column("t_int64", Int64),

		Column("t_float32", Float32),
		Column("t_float64", Float64),

		Column("t_string", String),
		//Column("t_binary", Binary),
		Column("t_rune", Rune),

		Column("t_bool", Bool),
	)

	def := Table("default_value", metadata,
		Column("id", Int).PrimaryKey(),

		Column("d_bool", Bool).Default(false),
		Column("d_int8", Int8).Default(int8(55)),
		Column("d_float32", Float32).Default(float32(10.2)),
		Column("d_string", String),
		//Column("d_binary", Binary).Default([]byte("123")),
		Column("d_rune", Rune).Default('a'),
	)

	// == Insert values
	types.InsertHelp("en",
		"int", "integer 8", "integer 16", "integer 32", "integer 64",
		"float 32", "float 64",
		"string", "rune",
		"boolean",
	)
	types.Insert(1, 8, 16, 32, 64, 1.32, 1.64, "one", "Z", true)

	def.InsertHelp("en", "id", "boolean", "integer 8", "float 32", "string", "rune")
	def.Insert(1, false, 10, 10.10, "foo", "z")
	// ==

	metadata.Create().Write()
}
