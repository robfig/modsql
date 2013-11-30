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

	//"github.com/kless/validate"
)

type constraintType int

const (
	primaryKey constraintType = 1 << iota
	foreignKey
	uniqueCons
)

type indexType int

const (
	_ indexType = iota
	noUniqIndex
	uniqIndex
)

// For columns with a wrong type
var columnsErr []string

type column struct {
	cons  constraintType
	index indexType

	type_ sqlType
	Name  string

	// Foreign key
	fkTable  string
	fkColumn string

	defaultValue interface{}
	//validators   validationType
}

// Column defines a new column.
func Column(name string, t sqlType) *column {
	c := new(column)
	c.Name = name
	c.type_ = t
	return c
}

// PrimaryKey defines the column to primary key.
func (c *column) PrimaryKey() *column {
	if c.cons == uniqueCons {
		c.addErrorCons()
	}
	if c.index != 0 {
		c.addErrorIndex()
	}

	c.cons |= primaryKey
	return c
}

// ForeignKey defines the column to foreign key.
func (c *column) ForeignKey(table, column string) *column {
	if c.cons == uniqueCons {
		c.addErrorCons()
	}
	if c.index != 0 {
		c.addErrorIndex()
	}

	c.cons |= foreignKey
	c.fkTable = table
	c.fkColumn = column
	return c
}

// Unique defines the column to UNIQUE constraint.
func (c *column) Unique() *column {
	if c.cons == primaryKey || c.cons == foreignKey {
		c.addErrorCons()
	}
	if c.index != 0 {
		c.addErrorIndex()
	}

	c.cons = uniqueCons
	return c
}

// Index sets an index.
func (c *column) Index(unique bool) *column {
	if c.cons != 0 {
		c.addErrorIndex()
	}

	if unique {
		c.index = uniqIndex
	} else {
		c.index = noUniqIndex
	}
	return c
}

// Default sets a value by default.
func (c *column) Default(v interface{}) *column {
	switch c.type_ {
	// MySQL: BLOB and TEXT columns cannot be assigned a default value.
	case String, Binary:
		log.Fatalf("type of column in %q can not have a default value", c.Name)
	}
	c.defaultValue = v

	if ok := c.checkDefValue(); !ok {
		columnsErr = append(columnsErr, fmt.Sprintf("\n column %q with type %T",
			c.Name, c.defaultValue),
		)
	}
	return c
}

/*// Validate sets some validator to ckeck the value in the actual column.
func (c *column) Validate(valid ...*validate.Validate) *column {
	for _, v := range valid {
		c.validators |= v
	}
	return c
}*/

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

	case uint8: // for the alias byte
		if c.type_ != Byte {
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

	case string:
		if c.type_ != String {
			return false
		}
	case []byte:
		if c.type_ != Binary {
			return false
		}

	//case time.Duration:
	//if c.type_ != Duration {
	//return false
	//}
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
		fmt.Sprintf("\n column %q only can have set a primary key, foreign key or unique constraint",
			c.Name))
}

func (c *column) addErrorIndex() {
	columnsErr = append(columnsErr,
		fmt.Sprintf("\n column %q only can have set an index or a constraint",
			c.Name))
}
