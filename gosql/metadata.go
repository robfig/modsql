// Copyright 2010  The "GoSQL" Authors
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

package gosql

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// The mode parameter to the Metadata function is a set of flags (or 0).
const Help uint = iota // Create tables related to the help.

// Defines a collection of table definitions.
type metadata struct {
	mode          uint
	useInsert     bool
	useInsertHelp bool
	tables        []*table
	queries       []byte
	model         []byte
}

// Initializes the type metadata.
func Metadata() *metadata {
	return &metadata{}
}

// Sets mode.
func (m *metadata) Mode(mode uint) *metadata {
	m.mode = mode
	return m
}

// * * *

// Issues both CREATE statements and Go definitions for all tables.
func (m *metadata) CreateAll() *metadata {
	create := make([]string, 0, 0)
	model := make([]string, 0, 0)

	pop := func(sl []string) []string {
		_, sl = sl[len(sl)-1], sl[:len(sl)-1]
		return sl
	}

	create = append(create, "BEGIN TRANSACTION;\n")
	model = append(model, header+"\n\npackage _RENAME_\n")

	for _, table := range m.tables {
		createLang := make([]string, 0, 0)

		if m.mode == Help {
			createLang = append(createLang,
				fmt.Sprintf("\nCREATE TABLE _%s (id TEXT PRIMARY KEY,\n", table.name))
		}

		create = append(create, fmt.Sprintf("\nCREATE TABLE %s (", table.name))
		model = append(model, fmt.Sprintf("\ntype %s struct {\n", table.name))

		for i, col := range table.columns {
			var field, extra string

			// The first field could not be a primary key
			if i == 0 {
				if !col.isPrimaryKey {
					create = append(create, "\n    ")
				}
			} else {
				field = "    "
			}

			model = append(model, fmt.Sprintf("%s %s\n", col.name, col.type_.Go()))
			create = append(create, fmt.Sprintf("%s %s",
				field+col.name, strings.ToUpper(col.type_.String())))

			if col.isPrimaryKey {
				extra += " PRIMARY KEY"
			}
			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch col.defaultValue.(type) {
				case string:
					extra += fmt.Sprintf("'%s'", col.defaultValue)

				case bool:
					// SQLite has not boolean type
					var conversion int

					if col.defaultValue.(bool) {
						conversion = 1
					}
					extra += fmt.Sprintf("%d", conversion)

				default:
					extra += fmt.Sprintf("%v", col.defaultValue)
				}
			}

			create = append(create, extra)

			// Add table for translation of fields comments
			if m.mode == Help && col.name != "id" {
				createLang = append(createLang, "    "+col.name+" TEXT")
				createLang = append(createLang, ",\n")
			}

			// The last column
			if i+1 == len(table.columns) {
				create = append(create, ");\n")
				model = append(model, "}\n")

				if m.mode == Help {
					createLang = pop(createLang)
					createLang = append(createLang, ");\n")
					create = append(create, createLang...)
				}
			} else {
				create = append(create, ",\n")
			}
		}
	}
	create = append(create, "\nCOMMIT;\n")

	// === Insert
	if m.useInsertHelp {
		m.insert(&create, _INSERT_HELP)
	}
	if m.useInsert {
		m.insert(&create, _INSERT_DATA)
	}

	m.queries = []byte(strings.Join(create, ""))
	m.model = []byte(strings.Join(model, ""))
	return m
}

// Writes SQL statements to a file or standard output.
func (m *metadata) Write(out output) {
	if out == FILEOUT {
		m.WriteTo(_SQL_OUTPUT, _MODEL_OUTPUT)
	} else if out == STDOUT {
		fmt.Printf("%s\n* * *\n\n", m.queries)
		m.format(os.Stdout)
	}
}

// Writes SQL statements to given files.
func (m *metadata) WriteTo(sqlFile, goFile string) {
	if len(m.queries) == 0 {
		fatal("No tables created. Use CreateAll()")
	}

	err := ioutil.WriteFile(sqlFile, m.queries, 0644)
	if err != nil {
		fatal("Failed to write file: %s", err)
	}

	file, err := os.OpenFile(goFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fatal("Failed to write file: %s", err)
	}
	defer file.Close()

	m.format(file)
	return
}

//
// === Utility

const (
	_INSERT_HELP uint = iota
	_INSERT_DATA
)

// To format Go source code.
const (
	_PARSER_MODE  = parser.ParseComments
	_PRINTER_MODE = printer.TabIndent | printer.UseSpaces
	_TAB_WIDTH    = 8
)

// Creates SQL statements to insert values; they are finally added to the main
// vector.
func (m *metadata) insert(main *[]string, value uint) {
	if value != _INSERT_HELP && value != _INSERT_DATA {
		fatal("argument \"value\" not valid for \"metadata.insert\": %d", value)
	}

	var data [][]interface{}
	insert := make([]string, 0, 0)
	insert = append(insert, "BEGIN TRANSACTION;\n")

	for _, table := range m.tables {
		tableName := table.name

		if value == _INSERT_HELP {
			data = table.help
			tableName = "_" + tableName
		} else if value == _INSERT_DATA {
			data = table.data
		}

		if len(data) != 0 {
			var columns []string

			for _, col := range table.columns {
				columns = append(columns, col.name)
			}

			for _, v := range data {
				insert = append(insert, fmt.Sprintf("\nINSERT INTO %q (%s) VALUES(%s);",
					tableName,
					strings.Join(columns, ", "),
					strings.Join(toString(v), ", ")))
			}
			insert = append(insert, "\n")
		}
	}

	insert = append(insert, "\nCOMMIT;\n")
	*main = append(*main, insert...)
}

// Converts a vector of interfaces to array of strings.
func toString(v []interface{}) (a []string) {
	for _, val := range v {
		switch val.(type) {
		case int:
			a = append(a, strconv.Itoa(val.(int)))
		case float32:
			a = append(a, strconv.FormatFloat(float64(val.(float32)), 'g', -1, 32))
		case float64:
			a = append(a, strconv.FormatFloat(val.(float64), 'g', -1, 64))
		case string:
			a = append(a, fmt.Sprintf("'%s'", val.(string)))
		case []uint8:
			a = append(a, fmt.Sprintf("'%s'", val.([]uint8)))
		case bool:
			b := 0
			if val == true {
				b = 1
			}
			a = append(a, strconv.Itoa(b))
		}
	}
	return
}

// Formats the Go source code.
func (m *metadata) format(out io.Writer) {
	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", m.model, _PARSER_MODE)
	if err != nil {
		fatal("Failed to format Go code: %s", err)
	}

	err = (&printer.Config{_PRINTER_MODE, _TAB_WIDTH}).Fprint(out, fset, ast)
	if err != nil {
		fatal("Failed to format Go code: %s", err)
	}

	return
}
