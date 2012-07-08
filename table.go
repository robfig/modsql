// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package modsql

import (
	"strings"
)

type table struct {
	name    string
	columns []column
	help    [][]interface{}
	data    [][]interface{}
	meta    *metadata
}

// Table defines a new table.
func Table(name string, meta *metadata, col ...*column) *table {
	if anyColumnErr {
		fatalf("Wrong type for default value in table %q: %s",
			name, strings.Join(columnsErr, ", "))
	}

	_table := new(table)
	_table.name = name
	_table.meta = meta

	for _, v := range col {
		_table.columns = append(_table.columns, *v)
	}

	meta.tables = append(meta.tables, _table)
	return _table
}

// Insert generates SQL statements to insert values.
func (t *table) Insert(a ...interface{}) {
	if len(a) != len(t.columns) {
		fatalf("incorrect number of arguments for Insert in table %q:"+
			" have %d, want %d",
			t.name, len(a), len(t.columns))
	}

	vec := make([]interface{}, 0, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	t.data = append(t.data, vec)
	t.meta.useInsert = true
}

// InsertHelp generates SQL statements to insert values on the help table.
func (t *table) InsertHelp(a ...string) {
	if t.meta.mode != Help {
		fatalf("Metadata Help mode is unset")
	}

	if len(a) != len(t.columns) {
		fatalf("incorrect number of arguments for Insert in table %q:"+
			" have %d, want %d",
			t.name, len(a), len(t.columns))
	}

	vec := make([]interface{}, 0, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	t.help = append(t.help, vec)
	t.meta.useInsertHelp = true
}
