// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import "strconv"

// sqlType represents the SQL type.
type sqlType byte

// SQL types, to be set in Column.
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

	Duration // time.Time
	DateTime // time.Date()
)

// goString returns the type corresponding to Go.
func (t sqlType) goString() string {
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

// sqlString returns the type corresponding to the engine used.
func (t sqlType) sqlString(engine sqlEngine) string {
	switch engine {

	// http://dev.mysql.com/doc/refman/5.6/en/data-types.html
	// http://nicj.net/mysql-text-vs-varchar-performance/
	case MySQL:
		switch t {
		case Bool:
			return "BOOL"

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

	// http://www.postgresql.org/docs/9.2/static/datatype-numeric.html
	case PostgreSQL:
		switch t {
		case Bool:
			return "boolean"

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

	// http://www.sqlite.org/datatype3.html
	case SQLite:
		switch t {
		case Bool:
			return "BOOL"

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
	}

	panic("unreachable")
}

// * * *

// formatBool returns the literal value for a boolean according to the SQL engine.
func (md *metadata) formatBool(b bool) string {
	if md.engine == SQLite {
		value := 0
		if b == true {
			value = 1
		}
		return strconv.Itoa(value)
	}

	value := "FALSE"
	if b == true {
		value = "TRUE"
	}
	return value
}
