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
	name string
	meta *metadata

	columns []column
	data    [][]interface{}

	// Constraints and indexes to table level
	uniqueCons []string
	pkCons     []string
	fkCons     []fkConstraint
	index      []compoIndex
}

// Table defines a new table.
func Table(name string, meta *metadata, col ...*column) *table {
	if len(columnsErr) != 0 {
		log.Fatalf("wrong type for default value in table %q: %s",
			name, strings.Join(columnsErr, ", "))
	}

	t := new(table)
	t.name = name
	t.meta = meta

	for _, v := range col {
		t.columns = append(t.columns, *v)
	}
	meta.tables = append(meta.tables, t)

	return t
}

// Insert generates SQL statements to insert values.
func (t *table) Insert(a ...interface{}) {
	if len(a) != len(t.columns) {
		log.Fatalf("incorrect number of arguments for Insert in table %q: have %d, want %d",
			t.name, len(a), len(t.columns))
	}

	vec := make([]interface{}, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	t.data = append(t.data, vec)
	t.meta.useInsert = true
}

// ForeignKey creates explicit/composite foreign key constraint.
// The keys in the map are the columns of this table, and the values are the
// foreign columns of the given table.
func (t *table) ForeignKey(table string, columns map[string]string) {
	if table == t.name {
		log.Fatalf("table %q: ForeignKey(): given foreign table can not have "+
			"the same name than actual table", table)
	}

	// Check foreign table
	found := false
	var tableColumns []column
	for _, v := range t.meta.tables {
		if v.name == table {
			tableColumns = v.columns
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("table %q: ForeignKey(): foreign table %q does not exist",
			t.name, table)
	}

	var fk fkConstraint

	for k, v := range columns {
		fk.src = append(fk.src, k)
		fk.dst = append(fk.dst, v)
	}

	t.checkColumns("ForeignKey", fk.src)

	for _, c := range fk.dst {
		found = false

		for _, tc := range tableColumns {
			if tc.name == c {
				if tc.cons != primaryKey && tc.cons != unique {
					log.Fatalf("table %q: ForeignKey(): column %q in foreign "+
						"table %q has to be a PRIMARY KEY or UNIQUE constraint",
						t.name, c, table)
				}

				found = true
				break
			}
		}
		if !found {
			log.Fatalf("table %q: ForeignKey(): foreign table %q has not column %q",
				t.name, table, c)
		}
	}

	fk.table = table
	t.fkCons = append(t.fkCons, fk)
}

// PrimaryKey creates explicit/composite primary key constraint.
func (t *table) PrimaryKey(columns ...string) {
	t.checkColumns("PrimaryKey", columns)
	t.pkCons = columns
}

// Unique creates explicit/composite unique constraint.
func (t *table) Unique(columns ...string) {
	t.checkColumns("Unique", columns)
	t.uniqueCons = columns
}

// Index creates an index on a group of columns.
func (t *table) Index(unique bool, columns ...string) {
	t.checkColumns("Index", columns)
	t.index = append(t.index, compoIndex{unique, columns})
}

// * * *

func (t *table) checkColumns(funcName string, columns []string) {
	for _, c := range columns {
		found := false

		for _, tc := range t.columns {
			if tc.name == c {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("table %q: %s(): column %q does not exist", t.name, funcName, c)
		}
	}
}
