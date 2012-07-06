// Copyright 2010  The "go2sql" Authors
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

package go2sql

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// mode defines the modes to use in metadata.Mode.
type mode byte

const (
	// If Help is set, it is created tables related to help.
	Help mode = iota + 1
)

// metadata defines a collection of table definitions.
type metadata struct {
	mode          mode
	useInsert     bool
	useInsertHelp bool
	tables        []*table
	queries       []byte
	model         []byte
}

// NewMetadata returns a new metadata.
func NewMetadata() *metadata {
	return &metadata{}
}

// Mode sets the mode.
func (md *metadata) Mode(m mode) *metadata {
	md.mode = m
	return md
}

// * * *

// CreateAll generates both CREATE statements and Go definitions for all tables.
func (md *metadata) CreateAll() *metadata {
	create := make([]string, 0, 0)
	model := make([]string, 0, 0)

	pop := func(sl []string) []string {
		_, sl = sl[len(sl)-1], sl[:len(sl)-1]
		return sl
	}

	create = append(create, "BEGIN TRANSACTION;\n")
	model = append(model, fmt.Sprintf("%s\n\npackage %s\n", header, getPkgName()))

	for _, table := range md.tables {
		createLang := make([]string, 0, 0)

		if md.mode == Help {
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
			if md.mode == Help && col.name != "id" {
				createLang = append(createLang, "    "+col.name+" TEXT")
				createLang = append(createLang, ",\n")
			}

			// The last column
			if i+1 == len(table.columns) {
				create = append(create, ");\n")
				model = append(model, "}\n")

				if md.mode == Help {
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
	if md.useInsertHelp {
		md.insert(&create, _INSERT_HELP)
	}
	if md.useInsert {
		md.insert(&create, _INSERT_DATA)
	}

	md.queries = []byte(strings.Join(create, ""))
	md.model = []byte(strings.Join(model, ""))
	return md
}

// Print prints both SQL statements and Go model.
func (md *metadata) Print() {
	fmt.Printf("%s\n* * *\n\n", md.queries)
	md.format(os.Stdout)
}

// Write writes both SQL statements and Go model to files using names by default.
func (md *metadata) Write() {
	md.WriteTo(_SQL_FILE, _MODEL_FILE)
}

// WriteTo writes both SQL statements and Go model to given files.
func (md *metadata) WriteTo(sqlFile, goFile string) {
	if len(md.queries) == 0 {
		fatalf("No tables created. Use CreateAll()")
	}

	err := ioutil.WriteFile(sqlFile, md.queries, 0644)
	if err != nil {
		fatalf("Failed to write file: %s", err)
	}

	file, err := os.OpenFile(goFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fatalf("Failed to write file: %s", err)
	}
	defer file.Close()

	md.format(file)
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

// insert generates SQL statements to insert values; they are finally added to
// the slice main.
func (md *metadata) insert(main *[]string, value uint) {
	if value != _INSERT_HELP && value != _INSERT_DATA {
		fatalf("argument \"value\" not valid for \"metadata.insert\": %d", value)
	}

	var data [][]interface{}
	insert := make([]string, 0, 0)
	insert = append(insert, "BEGIN TRANSACTION;\n")

	for _, table := range md.tables {
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

// toString converts to slice of strings.
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

// format formats the Go source code.
func (md *metadata) format(out io.Writer) {
	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", md.model, _PARSER_MODE)
	if err != nil {
		fatalf("Failed to format Go code: %s", err)
	}

	err = (&printer.Config{_PRINTER_MODE, _TAB_WIDTH}).Fprint(out, fset, ast)
	if err != nil {
		fatalf("Failed to format Go code: %s", err)
	}

	return
}

// == Utility
// ==

// getPkgName returns the package name of the actual directory.
func getPkgName() string {
	wd, err := os.Getwd()
	if err != nil {
		return "main"
	}

	if files, err := filepath.Glob("*.go"); err == nil && len(files) != 0 {
		for _, srcDir := range strings.Split(build.Default.GOPATH, ":") {
			importPath, err := filepath.Rel(srcDir, wd)
			if err != nil {
				continue
			}

			pkg, err := build.Import(importPath, srcDir, 0)
			if err != nil {
				continue
			}
			return pkg.Name
		}
	}

	return filepath.Base(wd)
}
