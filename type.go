// Copyright 2010  The "go2sql" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package go2sql

// sqlType represents the SQL type.
type sqlType uint8

// SQL types to set in Column.
const (
	Integer sqlType = iota + 1
	Float
	Text
	//Blob
	Boolean
)

var (
	sqlType_str = map[sqlType]string{
		Integer: "Integer",
		Float:   "Float",
		Text:    "Text",
		//Blob:    "Blob",
		Boolean: "Boolean",
	}

	sqlType_goType = map[sqlType]string{
		Integer: "int",
		Float:   "float32",
		Text:    "string",
		//Blob:    "[]byte",
		Boolean: "bool",
	}
)

func (t sqlType) String() string {
	return sqlType_str[t]
}

func (t sqlType) Go() string {
	return sqlType_goType[t]
}
