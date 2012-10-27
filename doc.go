// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package ModSQL enables to use a Go model to define the database model and
generate its corresponding SQL language and Go types. It is not an ORM neither
it is not going to be it since an ORM creates an extra layer to the database
access. The API is based in SQLAlchemy's (http://www.sqlalchemy.org/).

ModSQL enables to create primary key, foreign key and unique constraints, and
indexes at both column and table level.

It generates the files SQL and Go at writing to the file system, but it also can
shows the generated output. The name for the generated files start with "zmodsql".

If it is used the type Int, then the SQL files will have variables delimited by
"{{" and "}}", which will be parsed by the function Load according to the
architecture where it is being run.

Examples

See in directory "testdata"; the file "example.go" is the model, "zmodsql.go" is
the generated Go code, and "z[engine].sql" for the SQL files generated for every
engine which was indicated in the model (function Metadata).

For testing into a SQL engine, there is to run:

   go test -v -tags postgresql|mysql|sqlite

See files "db-*_test.go" to know how databases were configured.

Unsupported

The null handling is very different in every SQL engine (http://www.sqlite.org/nulls.html),
so instead I prefer to add empty values according to the type (just like in Go).  
Although they are useful for fields related to time and date.

Avoid cascades due to being magic; instead, I handle it from the application layer.  
http://stackoverflow.com/questions/59297/when-why-to-use-cascading-in-sql-server

Note

There are public methods which are not showed in the documentation due they
belong to private types. It happens in both types "column" and "table".
*/
package modsql
