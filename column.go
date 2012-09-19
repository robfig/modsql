// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"fmt"
	"log"
	"time"
)

// For columns with a wrong type
var (
	anyColumnErr bool
	columnsErr   []string
)

type column struct {
	isPrimaryKey bool
	isForeignKey bool

	type_ sqlType
	name  string

	// foreign key
	fkTable  string
	fkColumn string

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
	// MySQL: BLOB and TEXT columns cannot be assigned a default value.
	switch c.type_ {
	case String, Binary:
		log.Fatalf("type of column in %q can not have a default value", c.name)
	}

	c.defaultValue = v
	if ok := c.checkDefValue(); !ok {
		columnsErr = append(columnsErr, fmt.Sprintf("\n column %q with type %T",
			c.name, c.defaultValue),
		)
		anyColumnErr = true
	}
	return c
}

// PrimaryKey indicates that the column is a primary key.
func (c *column) PrimaryKey() *column {
	if c.isForeignKey {
		c.addErrorKey()
	}
	c.isPrimaryKey = true
	return c
}

// ForeignKey indicates that the column is foreign key.
func (c *column) ForeignKey(table, column string) *column {
	if c.isPrimaryKey {
		c.addErrorKey()
	}
	c.isForeignKey = true

	c.fkTable = table
	c.fkColumn = column
	return c
}

// * * *

// checkDefValue checks whether the default value has the correct type.
func (c *column) checkDefValue() bool {
	switch t := c.defaultValue.(type) {
	case bool:
		if c.type_ != Bool {
			return false
		}

	case int8:
		if c.type_ != Int8 {
			return false
		}
	case int16:
		if c.type_ != Int16 {
			return false
		}
	case int32: // rune is an alias
		if c.type_ != Int32 && c.type_ != Rune {
			return false
		}
	case int64:
		if c.type_ != Int64 {
			return false
		}

	case float32:
		if c.type_ != Float32 {
			return false
		}
	case float64:
		if c.type_ != Float64 {
			return false
		}

	case uint8: // for the alias byte
		if c.type_ != Byte {
			return false
		}

	case string:
		if c.type_ != String {
			return false
		}

	case []byte:
		if c.type_ != Binary {
			return false
		}

	case time.Duration:
		if c.type_ != Duration {
			return false
		}
	case time.Time:
		if c.type_ != DateTime {
			return false
		}

	default:
		panic(fmt.Sprintf("type %v not supported", t))
	}

	return true
}

func (c *column) addErrorKey() {
	columnsErr = append(columnsErr,
		fmt.Sprintf("\n column %q has several keys", c.name))
}
