// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package modsql allows use a Go model to define the database model and generate
its corresponding SQL language and Go types. It is not an ORM neither it is not
going to be it since an ORM creates an extra layer to the database access.
The API is based in SQLAlchemy's (http://www.sqlalchemy.org/).

It generates the files "model.sql" and "model.go" at writing to the file system;
it also can shows the generated output.

If it is used the type Int, then the SQL files will have variables delimited by
"{{" and "}}", which will be parsed by the function Load according to the
architecture where it is being run.

The function NewBuffer has the method Mode to create tables related to
localization. If it is set, then at generating SQL it is created an extra table
(starting wich "_") for each model.

NOTE: it is not ready to working with relations between tables since I don't
need it by now.


Operating instructions

This example is used in file "modsql_test.go" except that it writes to file.
To run it, use "go run file.go".

	package main

	import . "github.com/kless/modsql"

	func main() {
		metadata := NewMetadata(PostgreSQL).Mode(Help)

		types := Table("types", metadata,
			Column("id", Int64).PrimaryKey(),
			Column("t_int", Int32),
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

		metadata.CreateAll().Write()
	}
*/
package modsql
