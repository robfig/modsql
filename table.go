// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"log"
	"strings"
)

// A ForeignColumn represents the relationship between columns of two tables.
type ForeignColumn struct {
	Local   string
	Foreign string
}

type fkConstraint struct {
	table string
	src   []string
	dst   []string
}

type compoIndex struct {
	isUnique bool
	index    []string
}

type table struct {
	isEnum    bool // table with list of permitted values that are enumerated
	startEnum int

	Name    string
	sqlName string
	meta    *metadata
	Columns []column

	// Constraints and indexes to table level
	uniqueCons []string
	pkCons     []string
	fkCons     []fkConstraint
	index      []compoIndex

	// To insert values
	data     [][]interface{}
	testData [][]interface{}
}

// Enum defines a table for enumeration values starting from start.
func Enum(name string, meta *metadata, intType sqlType, start int, value ...string) {
	if intType < Int || intType > Int64 {
		log.Fatalf("wrong type for integer: %s", intType.goString())
	}

	t := Table(name, meta,
		Column("id", intType).PrimaryKey(),
		Column("name", String),
	)
	t.isEnum = true
	t.startEnum = start

	for i, v := range value {
		t.Insert(i+start, v)
	}
}

// Table defines a new table.
func Table(name string, meta *metadata, col ...*column) *table {
	if len(columnsErr) != 0 {
		log.Fatalf("wrong type for default value in table %q: %s",
			name, strings.Join(columnsErr, ", "))
	}

	t := new(table)
	t.Name = name
	t.sqlName = quoteSQL(name)
	t.meta = meta

	for _, v := range col {
		t.Columns = append(t.Columns, *v)
	}
	meta.tables = append(meta.tables, t)

	return t
}

// Insert generates SQL statements to insert values.
func (t *table) Insert(a ...interface{}) {
	if len(a) != len(t.Columns) {
		log.Fatalf("incorrect number of arguments to insert in table %q: have %d, want %d",
			t.Name, len(a), len(t.Columns))
	}

	vec := make([]interface{}, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	t.data = append(t.data, vec)
	t.meta.useInsert = true
}

// InsertTestData generates SQL statements to insert values in test database.
// It is generated in file names with suffix "_test".
func (t *table) InsertTestData(a ...interface{}) {
	if len(a) != len(t.Columns) {
		log.Fatalf("incorrect number of arguments to insert test data in table %q: have %d, want %d",
			t.Name, len(a), len(t.Columns))
	}

	vec := make([]interface{}, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	t.testData = append(t.testData, vec)
	t.meta.useInsertTest = true
}

// ForeignKey creates explicit/composite foreign key constraint.
// The keys in the map are the columns of this table, and the values are the
// foreign columns of the given table.
func (t *table) ForeignKey(table string, columns ...ForeignColumn) {
	if table == t.Name {
		log.Fatalf("table %q: ForeignKey(): given foreign table can not have "+
			"the same name than actual table", table)
	}

	// Check foreign table
	found := false
	var tableColumns []column
	for _, v := range t.meta.tables {
		if v.Name == table {
			tableColumns = v.Columns
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("table %q: ForeignKey(): foreign table %q does not exist",
			t.Name, table)
	}

	var fk fkConstraint

	for _, col := range columns {
		fk.src = append(fk.src, col.Local)
		fk.dst = append(fk.dst, col.Foreign)
	}

	t.existColumns("ForeignKey", fk.src)

	for _, c := range fk.dst {
		found = false

		for _, tc := range tableColumns {
			if tc.Name == c {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("table %q: ForeignKey(): foreign table %q has not column %q",
				t.Name, table, c)
		}
	}

	fk.table = table
	t.fkCons = append(t.fkCons, fk)
}

// PrimaryKey creates explicit/composite primary key constraint.
func (t *table) PrimaryKey(columns ...string) {
	t.existColumns("PrimaryKey", columns)
	t.pkCons = columns
}

// Index creates an index on a group of columns.
func (t *table) Index(unique bool, columns ...string) {
	t.existColumns("Index", columns)
	t.index = append(t.index, compoIndex{unique, columns})
}

// Unique creates explicit/composite unique constraint.
func (t *table) Unique(columns ...string) {
	t.existColumns("Unique", columns)
	t.uniqueCons = columns
}

// * * *

// existColumns checks if the given columns are in the actual table.
func (t *table) existColumns(funcName string, columns []string) {
	for _, c := range columns {
		found := false

		for _, tc := range t.Columns {
			if tc.Name == c {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("table %q: %s(): column %q does not exist", t.Name, funcName, c)
		}
	}
}
