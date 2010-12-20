// Copyright 2010  The "SQLModel" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package sqlmodel

import (
	"container/vector"
	"strings"
)

type table struct {
	name    string
	columns []column
	values  []*vector.Vector
}


func Table(name string, meta *metadata, col ...*column) *table {
	if anyColumnErr {
		fatal("Wrong type for default value in table %q: %s",
			name, strings.Join(columnsErr, ", "))
	}

	_table := new(table)
	_table.name = name

	for _, v := range col {
		_table.columns = append(_table.columns, *v)
	}

	meta.tables = append(meta.tables, _table)
	return _table
}

// Creates SQL statements to insert values.
func (self *table) Insert(a... interface{}) {
	if len(a) != len(self.columns) {
		fatal("incorrect number of arguments for Insert in table %q:"+
			" have %d, want %d", self.name, len(a), len(self.columns))
	}

	vec := new(vector.Vector)
	for _, v := range a {
		vec.Push(v)
	}

	self.values = append(self.values, vec)
}

