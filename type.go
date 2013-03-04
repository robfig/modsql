// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import "fmt"

// An Engine represents the SQL engine.
type Engine int

const (
	MySQL Engine = iota + 1
	Postgres
	SQLite
)

func (e Engine) String() string {
	switch e {
	case MySQL:
		return "MySQL"
	case Postgres:
		return "Postgres"
	case SQLite:
		return "SQLite"
	}
	panic("unreachable")
}

func (e Engine) check() error {
	switch e {
	case MySQL, Postgres, SQLite:
		return nil
	}
	return fmt.Errorf("wrong engine: %s", e)
}

// quoteChar are the characters used to quote a name according to a SQL engine.
var quoteChar = map[Engine]string{
	MySQL:    "`",
	Postgres: `"`,
	SQLite:   `"`,
}

// * * *

// sqlType represents the SQL type.
type sqlType byte

// SQL types, to be set in Column.
const (
	Bool sqlType = iota + 1

	Int
	Int8
	Int16
	Int32
	Int64

	Byte
	Rune

	Float32
	Float64

	String
	Binary

	//Duration // time.Duration
	DateTime // time.Time
)

// goString returns the type corresponding to Go.
func (t sqlType) goString() string {
	switch t {
	case Bool:
		return "bool"

	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"

	case Byte:
		return "byte" // uint8 (fits into int16)
	case Rune:
		return "rune" // int32

	case Float32:
		return "float32"
	case Float64:
		return "float64"

	case String:
		return "string"
	case Binary:
		return "[]byte"

	//case Duration:
		//return "time.Duration"
	case DateTime:
		return "time.Time"
	}

	panic("unreachable")
}

// tmplAction returns a template action which will enable to generate the SQL type
// for every SQL engine.
func (t sqlType) tmplAction() string {
	switch t {
	case Bool:
		return "{{.Bool}}"

	case Int:
		return "{{.Int}}"
	case Int8:
		return "{{.Int8}}"
	case Int16:
		return "{{.Int16}}"
	case Int32:
		return "{{.Int32}}"
	case Int64:
		return "{{.Int64}}"

	case Byte:
		return "{{.Byte}}"
	case Rune:
		return "{{.Rune}}"

	case Float32:
		return "{{.Float32}}"
	case Float64:
		return "{{.Float64}}"

	case String:
		return "{{.String}}"
	case Binary:
		return "{{.Binary}}"

	//case Duration:
		//return "{{.Duration}}"
	case DateTime:
		return "{{.DateTime}}"
	}

	panic("unreachable")
}

// boolAction returns the template action for a boolean.
func boolAction(b bool) string {
	if b == true {
		return "{{.True}}"
	}
	return "{{.False}}"
}

// * * *

// A sqlAction represents data to pass to the SQL template.
type sqlAction struct {
	Engine string

	Bool  string
	True  string
	False string

	Int   string
	Int8  string
	Int16 string
	Int32 string
	Int64 string

	Byte string
	Rune string

	Float32 string
	Float64 string

	String      string
	StringLimit string
	Binary      string

	//Duration string
	DateTime string

	Q string // character of quote

	MySQLDrop0   string
	MySQLDrop1   string
	PostgresDrop string
}

// getSQLAction returns data corresponding to the engine used.
func getSQLAction(eng Engine) *sqlAction {
	a := new(sqlAction)

	switch eng {

	// http://dev.mysql.com/doc/refman/5.6/en/data-types.html
	// http://nicj.net/mysql-text-vs-varchar-performance/
	case MySQL:
		a = &sqlAction{
			Bool: "BOOL",

			Int:   "{{.MySQLInt}}", // to be parsed in function Load
			Int8:  "TINYINT",
			Int16: "SMALLINT",
			Int32: "INT",
			Int64: "BIGINT",

			Byte: "SMALLINT",
			Rune: "INT",

			Float32: "FLOAT",
			Float64: "DOUBLE",

			String:      "TEXT",
			StringLimit: "VARCHAR(255)",
			Binary:      "BLOB",

			//Duration: "TIME",
			DateTime: "TIMESTAMP",

			Q: quoteChar[MySQL],

			MySQLDrop0: "\nSET FOREIGN_KEY_CHECKS=0;\n",
			MySQLDrop1: "\n\nSET FOREIGN_KEY_CHECKS=1;",
		}

	// http://www.postgresql.org/docs/9.2/static/datatype-numeric.html
	case Postgres:
		a = &sqlAction{
			Bool: "boolean",

			Int:   "{{.PostgresInt}}", // to be parsed in function Load
			Int8:  "smallint",
			Int16: "smallint",
			Int32: "integer",
			Int64: "bigint",

			Byte: "smallint",
			Rune: "integer",

			Float32: "real",
			Float64: "double precision",

			String:      "text",
			StringLimit: "text",
			Binary:      "bytea",

			//Duration: "time without time zone",
			DateTime: "timestamp with time zone",

			Q: quoteChar[Postgres],

			PostgresDrop: " CASCADE", // automatically drop objects that depend on the table
		}

	// http://www.sqlite.org/datatype3.html
	case SQLite:
		a = &sqlAction{
			Bool: "BOOL",

			Int:   "INTEGER",
			Int8:  "INTEGER",
			Int16: "INTEGER",
			Int32: "INTEGER",
			Int64: "INTEGER",

			Byte: "INTEGER",
			Rune: "INTEGER",

			Float32: "REAL",
			Float64: "REAL",

			String:      "TEXT",
			StringLimit: "TEXT",
			Binary:      "BLOB",

			//Duration: "INTEGER", // time()
			DateTime: "TEXT",    // datetime()

			Q: quoteChar[SQLite],
		}
	}

	if eng == SQLite {
		a.False = "0"
		a.True = "1"
	} else {
		a.False = "FALSE"
		a.True = "TRUE"
	}
	a.Engine = eng.String()

	return a
}
