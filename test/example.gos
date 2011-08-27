#!/usr/bin/goscript

package main

import . "github.com/kless/SQLModel/sqlmodel"

func main() {
	metadata := Metadata().Mode(Help)

	types := Table("types", metadata,
		Column("id", Integer).PrimaryKey(),
		Column("t_int", Integer),
		Column("t_float", Float),
		Column("t_text", Text),
		Column("t_blob", Blob),
		Column("t_bool", Boolean),
	)

	def := Table("default_value", metadata,
		Column("id", Integer).PrimaryKey(),
		Column("d_int", Integer).Default(55),
		Column("d_float", Float).Default(10.2),
		Column("d_text", Text).Default("string"),
		//Column("d_blob", Blob).Default([]byte("123")),
		Column("d_bool", Boolean).Default(false),
	)

	// === Insert values
	types.InsertHelp("en", "integer", "float", "text", "binary", "boolean")
	types.Insert(1, 10, 1.1, "one", []byte("one"), true)
	types.Insert(2, 20, 2.2, "two", []byte("two"), false)

	def.InsertHelp("en", "integer", "float", "text", "boolean")
	def.Insert(1, 10, 10.1, "foo", true)
	// ===

	metadata.CreateAll().Write(STDOUT)
	//metadata.Write(FILEOUT)
}

