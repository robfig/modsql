// Copyright 2010  The "go2sql" Authors
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

/*
Package go2sql allows use a Go to define a database model and generate its
corresponding SQL language and Go types. It is not an ORM neither it is not
goint to be it since an ORM creates an extra layer to the database access.

It generates the files "model.sql" and "model.go" at writing to the file system.
But it also can show the generated output.

The API is based in SQLAlchemy's (http://www.sqlalchemy.org/) and here it is all
all SQLAlchmey's types (http://www.sqlalchemy.org/docs/core/types.html),
althought I have only added some basic types.

The function Metadata has method Mode to create tables related to localization.
If it is set, then at generating SQL, it creates an extra table (starting wich
"_") for each model.

NOTE: it is tested with SQLite3, and it is not ready to working with relations
between tables since I don't need it by now.


Operating instructions

Here it is the example used in "go2sql_test.go" except that this one writes to
file. To run it, use "go run file.go".

	package main

	import . "github.com/kless/go2sql"

	func main() {
		metadata := NewMetadata().Mode(Help)

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

		// == Insert values
		types.InsertHelp("en", "integer", "float", "text", "binary", "boolean")
		types.Insert(1, 10, 1.1, "one", []byte("one"), true)
		types.Insert(2, 20, 2.2, "two", []byte("two"), false)

		def.InsertHelp("en", "integer", "float", "text", "boolean")
		def.Insert(1, 10, 10.1, "foo", true)
		// ==

		metadata.CreateAll().Write()
	}
*/
package go2sql

import (
	"fmt"
	"os"
)

const (
	_SQL_FILE   = "model.sql"
	_MODEL_FILE = "model.go" // Go definitions related to each SQL table

	header = "// MACHINE GENERATED BY \"github.com/kless/go2sql\"\n// ==="
)

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
