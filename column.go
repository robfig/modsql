// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package modsql

// For columns with a wrong type
var (
	anyColumnErr bool
	columnsErr   []string
)

type column struct {
	name         string
	type_        sqlType
	isPrimaryKey bool
	defaultValue interface{}
}

// Column defines a new column.
func Column(name string, type_ sqlType) *column {
	col := new(column)
	col.name = name
	col.type_ = type_
	return col
}

// Default sets a value by default.
func (c *column) Default(v interface{}) *column {
	c.defaultValue = v

	if ok := c.check(); !ok {
		columnsErr = append(columnsErr, c.name)
		anyColumnErr = true
	}

	return c
}

// PrimaryKey indicates that the column is a primary key.
func (c *column) PrimaryKey() *column {
	c.isPrimaryKey = true
	return c
}

// check checks whether the value by default has the correct type.
func (c *column) check() bool {
	switch c.defaultValue.(type) {
	case bool:
		if c.type_ != Boolean {
			return false
		}
	case float32, float64:
		if c.type_ != Float {
			return false
		}
	case int:
		if c.type_ != Integer {
			return false
		}
	case string:
		if c.type_ != Text {
			return false
		}
	default:
		return false
	}

	return true
}
