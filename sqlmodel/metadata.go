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
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// To format Go source code.
const (
	_PARSER_MODE  = parser.ParseComments
	_PRINTER_MODE = printer.TabIndent | printer.UseSpaces
	_TAB_WIDTH    = 8
)

// Defines a collection of table definitions.
type metadata struct {
	tables  []*table
	queries []byte
	model   []byte
}


func Metadata() *metadata {
	return &metadata{}
}

// Issues both CREATE statements and Go definitions for all tables.
func (self *metadata) CreateAll() *metadata {
	var create, model vector.StringVector

	create.Push("BEGIN TRANSACTION;\n")
	model.Push("// MACHINE GENERATED.\n\npackage _\n")

	for _, table := range self.tables {
		var createLang vector.StringVector
		createLang.Push(fmt.Sprintf("\nCREATE TABLE _%s (id TEXT PRIMARY KEY,\n",
			table.name))

		create.Push(fmt.Sprintf("\nCREATE TABLE %s (", table.name))
		model.Push(fmt.Sprintf("\ntype %s struct {\n", table.name))

		for i, col := range table.columns {
			var field, extra string

			// The first field could not be a primary key
			if i == 0 {
				if !col.isPrimaryKey {
					create.Push("\n    ")
				}
			} else {
				field = "    "
			}

			model.Push(fmt.Sprintf("%s %s\n", col.name, col.type_.Go()))
			create.Push(fmt.Sprintf("%s %s",
				field+col.name, strings.ToUpper(col.type_.String())))

			if col.isPrimaryKey {
				extra += " PRIMARY KEY"
			}
			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch col.defaultValue.(type) {
				case string:
					extra += fmt.Sprintf("%q", col.defaultValue)

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

			create.Push(extra)

			// Add table for translation of fields comments
			if col.name != "id" {
				createLang.Push("    " + col.name + " TEXT")
				createLang.Push(",\n")
			}

			// The last column
			if i+1 == len(table.columns) {
				create.Push(");\n")
				model.Push("}\n")

				createLang.Pop()
				createLang.Push(");\n")
				create.AppendVector(&createLang)
			} else {
				create.Push(",\n")
			}
		}
	}

	create.Push("\nCOMMIT;\n")

	self.queries = []byte(strings.Join(create, ""))
	self.model = []byte(strings.Join(model, ""))
	return self
}

// Writes SQL statements to a file or standard output.
func (self *metadata) Write(out output) {
	if out == FILEOUT {
		self.WriteTo(_SQL_OUTPUT, _MODEL_OUTPUT)
	} else if out == STDOUT {
		fmt.Printf("%s\n* * *\n\n", self.queries)
		self.format(os.Stdout)
	}
}

// Writes SQL statements to given files.
func (self *metadata) WriteTo(sqlFile, goFile string) {
	if len(self.queries) == 0 {
		fatal("No tables created. Use CreateAll()")
	}

	err := ioutil.WriteFile(sqlFile, self.queries, 0644)
	if err != nil {
		goto _error
	}

	file, err := os.Open(goFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		goto _error
	}
	defer file.Close()

	self.format(file)
	return

_error:
	fatal("Failed to write file: %s", err)
}

// Formats the Go source code.
func (self *metadata) format(out io.Writer) {
	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", self.model, _PARSER_MODE)
	if err != nil {
		goto _error
	}

	_, err = (&printer.Config{_PRINTER_MODE, _TAB_WIDTH, nil}).Fprint(out, fset, ast)
	if err != nil {
		goto _error
	}

	return

_error:
	fatal("Failed to format Go source code: %s", err)
}

