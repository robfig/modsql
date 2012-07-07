/*
Package go2sql allows use a Go to define a database model and generate its
corresponding SQL language and Go types. It is not an ORM neither it is not
goint to be it since an ORM creates an extra layer to the database access.

It generates the files "model.sql" and "model.go" at writing to the file system.
But it also can show the generated output.

The API is based in SQLAlchemy's (http://www.sqlalchemy.org/), althought it has
been only added some basic types.

The function NewBuffer has the method Mode to create tables related to
localization. If it is set, then at generating SQL it is created an extra table
(starting wich "_") for each model.

NOTE: it is tested with SQLite3 and PostgreSQL. And it is not ready to working
with relations between tables since I don't need it by now.


Operating instructions

Here it is the example used in "go2sql_test.go" except that this one writes to
file. To run it, use "go run file.go".

	package main

	import . "github.com/kless/go2sql"

	func main() {
		metadata := NewMetadata(PostgreSQL).Mode(Help)

		types := Table("types", metadata,
			Column("id", Integer).PrimaryKey(),
			Column("t_int", Integer),
			Column("t_float", Float),
			Column("t_text", Text),
			Column("t_bool", Boolean),
		)

		def := Table("default_value", metadata,
			Column("id", Integer).PrimaryKey(),
			Column("d_int", Integer).Default(55),
			Column("d_float", Float).Default(10.2),
			Column("d_text", Text).Default("string"),
			Column("d_bool", Boolean).Default(false),
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
package go2sql
