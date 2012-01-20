// Copyright 2010  The "GotoSQL" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tosql

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

func Table(name string, meta *metadata, col ...*column) *table {
	if anyColumnErr {
		fatal("Wrong type for default value in table %q: %s",
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

// Creates SQL statements to insert values.
func (self *table) Insert(a ...interface{}) {
	if len(a) != len(self.columns) {
		fatal("incorrect number of arguments for Insert in table %q:"+
			" have %d, want %d",
			self.name, len(a), len(self.columns))
	}

	vec := make([]interface{}, 0, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	self.data = append(self.data, vec)
	self.meta.useInsert = true
}

// Creates SQL statements to insert values on its help table.
func (self *table) InsertHelp(a ...string) {
	if self.meta.mode != Help {
		fatal("Metadata Help mode is unset")
	}

	if len(a) != len(self.columns) {
		fatal("incorrect number of arguments for Insert in table %q:"+
			" have %d, want %d",
			self.name, len(a), len(self.columns))
	}

	vec := make([]interface{}, 0, 0)
	for _, v := range a {
		vec = append(vec, v)
	}

	self.help = append(self.help, vec)
	self.meta.useInsertHelp = true
}
