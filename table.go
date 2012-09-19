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

type table struct {
	name string
	meta *metadata

	columns []column
	data    [][]interface{}
}

// Table defines a new table.
func Table(name string, meta *metadata, col ...*column) *table {
	if anyColumnErr {
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
