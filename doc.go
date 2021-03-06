// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package ModSQL defines the database model and generates its corresponding SQL
language and Go types. It is not an ORM neither it is not going to be it since
an ORM creates an extra layer to the database access. The API is based in
SQLAlchemy's (http://www.sqlalchemy.org/).

ModSQL enables to create primary key, foreign key and unique constraints, and
indexes at both column and table level.

It generates the SQL and Go files at writing to the file system, but it also can
shows the generated output.

If it is used the type Int, then the SQL files will have variables delimited by
"{{" and "}}", which will be parsed by the function Load according to the
architecture where it is being run.

Features

Dialect implemented for PostgreSQL, MySQL, SQLite3
Schema generation
Support primary and foreign keys, indexes and unique constraints, also for composites
Default values
Enumerations

Enumeration

The function "Enum" allows to create a table with the given names whose values
will be the same in both SQL tables and Go code.

Some SQL engines have a type to handle enumerations but they have some issues
as explained here:

http://komlenic.com/244/8-reasons-why-mysqls-enum-data-type-is-evil/

Datetime

That data must be stored in UTC. By this reason, the data type for DateTime in
PostgreSQL is defined with "timestamp without time zone".

It is used "time.Time{}" to get the initial value to zero, which is better than
using NULL values.

Unsupported

The null handling is very different in every SQL engine (http://www.sqlite.org/nulls.html),
so instead I prefer to add empty values according to the type (just like in Go).

time.Duration is not supported by sql.Scanner: code.google.com/p/go/issues/detail?id=4954

Examples

The directory 'test/data/sql' has the files generated from 'test/modeler.go'
which is run through 'gotask init'.

"[engine]*.sql" are the SQL files for every engine indicated in the model
(function Metadata in 'test/modeler.go').

For testing into a SQL engine, there is to run:

   test> gotask test-postgres|test-mysql|test-sqlite

See files 'test/[engine]_task.go' to know how databases were configured.

Avoid cascades due to being magic; instead, I handle it from the application layer.
http://stackoverflow.com/questions/59297/when-why-to-use-cascading-in-sql-server

Usage

You have to create a directory for the model's file or files; as suggestion,
name it "ModSQL". Then, from the project's directory run "go run ModSQL/[file].go"

The Go file generated uses the constant "ENGINE" according to the given values
in the function "Metadata", using the first engine by default.

Note

There are public methods which are not showed in the documentation due they
belong to private types. It happens in both types "column" and "table".
*/
package modsql
