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

type constraint int

const (
	cNone constraint = iota
	cPrimaryKey
	cForeignKey
	cUnique
)

type column struct {
	cons constraint

	type_ sqlType
	name  string

	// Foreign key
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

// ForeignKey defines the column to foreign key.
func (c *column) ForeignKey(table, column string) *column {
	if c.cons == cPrimaryKey || c.cons == cUnique {
		c.addErrorCons()
	}
	c.cons = cForeignKey
	c.fkTable = table
	c.fkColumn = column
	return c
}

// PrimaryKey defines the column to primary key.
func (c *column) PrimaryKey() *column {
	if c.cons == cForeignKey || c.cons == cUnique {
		c.addErrorCons()
	}
	c.cons = cPrimaryKey
	return c
}

// Unique defines the column to constraint unique.
func (c *column) Unique() *column {
	if c.cons == cPrimaryKey || c.cons == cForeignKey {
		c.addErrorCons()
	}
	c.cons = cUnique
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

func (c *column) addErrorCons() {
	columnsErr = append(columnsErr,
		fmt.Sprintf("\n column %q only can be primary key, foreign key or unique",
			c.name))
}
