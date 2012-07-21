// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

// sqlType represents the SQL type.
type sqlType byte

// SQL types to set in Column.
const (
	Bool sqlType = iota + 1

	//Int
	Int8
	Int16
	Int32
	Int64

	Float32
	Float64

	String
	Byte
	Rune

	Binary

	Duration
	DateTime // time.Date()
)

// GoString returns the type corresponding to Go.
func (t sqlType) GoString() string {
	switch t {
	case Bool:
		return "bool"

	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"

	case Float32:
		return "float32"
	case Float64:
		return "float64"

	case String:
		return "string"
	case Byte:
		return "byte" // uint8
	case Rune:
		return "rune" // int32

	case Binary:
		return "[]byte"

	case Duration:
		return "*time.Duration"
	case DateTime:
		return "time.Time"
	}

	panic("unreachable")
}

// SQLString returns the type corresponding to the engine used.
func (t sqlType) SQLString(engine sqlEngine) string {
	switch engine {
	case MySQL:
		return t.mysql()
	case PostgreSQL:
		return t.postgre()
	case SQLite:
		return t.sqlite()
	}
	panic("unreachable")
}

// http://dev.mysql.com/doc/refman/5.6/en/data-types.html
// http://nicj.net/mysql-text-vs-varchar-performance/
//
// MySQL

// mysql returns the data type corresponding to MySQL.
func (t sqlType) mysql() string {
	switch t {
	case Bool:
		return "BOOL" // TRUE, FALSE

	case Int8:
		return "TINYINT"
	case Int16:
		return "SMALLINT"
	case Int32:
		return "INT"
	case Int64:
		return "BIGINT"

	case Float32:
		return "FLOAT"
	case Float64:
		return "DOUBLE"

	case String:
		return "TEXT"
	case Byte:
		return "CHAR(1)"
	case Rune:
		return "CHAR(4)"

	case Binary:
		return "BLOB"

	case Duration:
		return "TIME"
	case DateTime:
		return "TIMESTAMP"
	}

	panic("unreachable")
}

// http://www.postgresql.org/docs/9.2/static/datatype-numeric.html
//
// PostgreSQL

// postgre returns the data type corresponding to PostgreSQL.
func (t sqlType) postgre() string {
	switch t {
	case Bool:
		return "boolean" // TRUE, FALSE

	case Int8, Int16:
		return "smallint"
	case Int32:
		return "integer"
	case Int64:
		return "bigint"

	case Float32:
		return "real"
	case Float64:
		return "double precision"

	case String:
		return "text"
	case Byte:
		return "character"
	case Rune:
		return "character varying(4)"

	case Binary:
		return "bytea"

	case Duration:
		return "time without time zone"
	case DateTime:
		return "timestamp with time zone"
	}

	panic("unreachable")
}

// http://www.sqlite.org/datatype3.html
//
// SQLite

// sqlite returns the data type corresponding to SQLite.
func (t sqlType) sqlite() string {
	switch t {
	case Bool:
		return "BOOL" // 0, 1

	case Int8, Int16, Int32, Int64:
		return "INTEGER"

	case Float32, Float64:
		return "REAL"

	case String, Byte, Rune:
		return "TEXT"

	case Binary:
		return "BLOB"

	case Duration:
		return "INTEGER" // time()
	case DateTime:
		return "TEXT" // datetime()
	}

	panic("unreachable")
}
