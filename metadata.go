// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package modsql

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

// mode represents the modes to use in metadata.Mode.
type mode byte

const (
	// If Help is set, it is created tables related to help.
	Help mode = iota
)

// dbEngine represents the SQL engine.
type sqlEngine byte

// SQL engines.
const (
	SQLite sqlEngine = iota
	MySQL
	PostgreSQL
)

// metadata defines a collection of table definitions.
type metadata struct {
	engine        sqlEngine
	mode          mode
	useInsert     bool
	useInsertHelp bool
	tables        []*table
	queries       []byte
	model         []byte
}

// NewMetadata returns a new metadata.
func NewMetadata(engine sqlEngine) *metadata {
	return &metadata{engine: engine}
}

// Mode sets the mode.
func (md *metadata) Mode(m mode) *metadata {
	md.mode = m
	return md
}

// * * *

// CreateAll generates both SQL statements and Go definitions for all tables.
func (md *metadata) CreateAll() *metadata {
	sql := make([]string, 0, 0)
	model := make([]string, 0, 0)

	pop := func(sl []string) []string {
		_, sl = sl[len(sl)-1], sl[:len(sl)-1]
		return sl
	}

	sql = append(sql, fmt.Sprintf("%s\nBEGIN TRANSACTION;\n", header))
	model = append(model, fmt.Sprintf("%s\npackage %s\n", header, getPkgName()))

	for _, table := range md.tables {
		sqlLang := make([]string, 0, 0)

		if md.mode == Help {
			sqlLang = append(sqlLang,
				fmt.Sprintf("\nCREATE TABLE _%s (id TEXT PRIMARY KEY,\n", table.name))
		}

		sql = append(sql, fmt.Sprintf("\nCREATE TABLE %s (", table.name))
		model = append(model, fmt.Sprintf("\ntype %s struct {\n", table.name))

		for i, col := range table.columns {
			var field, extra string

			// The first field could not be a primary key
			if i == 0 {
				if !col.isPrimaryKey {
					sql = append(sql, "\n    ")
				}
			} else {
				field = "    "
			}

			model = append(model, fmt.Sprintf("%s %s\n", col.name, col.type_.Go()))
			sql = append(sql, fmt.Sprintf("%s %s",
				field+col.name, strings.ToUpper(col.type_.String())))

			if col.isPrimaryKey {
				extra += " PRIMARY KEY"
			}
			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch t := col.defaultValue.(type) {
				case string:
					extra += fmt.Sprintf("'%s'", t)
				case bool:
					extra += fmt.Sprintf("%s", md.getbool(t))
				default:
					extra += fmt.Sprintf("%v", t)
				}
			}

			sql = append(sql, extra)

			// Add table for translation of fields comments
			if md.mode == Help && col.name != "id" {
				sqlLang = append(sqlLang, "    "+col.name+" TEXT")
				sqlLang = append(sqlLang, ",\n")
			}

			// The last column
			if i+1 == len(table.columns) {
				sql = append(sql, ");\n")
				model = append(model, "}\n")

				if md.mode == Help {
					sqlLang = pop(sqlLang)
					sqlLang = append(sqlLang, ");\n")
					sql = append(sql, sqlLang...)
				}
			} else {
				sql = append(sql, ",\n")
			}
		}
	}
	sql = append(sql, "\nCOMMIT;\n")

	// == Insert
	if md.useInsertHelp {
		md.insert(&sql, _INSERT_HELP)
	}
	if md.useInsert {
		md.insert(&sql, _INSERT_DATA)
	}

	md.queries = []byte(strings.Join(sql, ""))
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

// * * *

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
					strings.Join(md.toString(v), ", ")))
			}
			insert = append(insert, "\n")
		}
	}

	insert = append(insert, "\nCOMMIT;\n")
	*main = append(*main, insert...)
}

// toString converts to slice of strings.
func (md *metadata) toString(v []interface{}) []string {
	res := make([]string, 0)

	for _, val := range v {
		switch t := val.(type) {
		case int:
			res = append(res, strconv.Itoa(t))
		case float32:
			res = append(res, strconv.FormatFloat(float64(t), 'g', -1, 32))
		case float64:
			res = append(res, strconv.FormatFloat(t, 'g', -1, 64))
		case string:
			res = append(res, fmt.Sprintf("'%s'", t))
		case []uint8:
			res = append(res, fmt.Sprintf("'%s'", t))
		case bool:
			res = append(res, md.getbool(t))
		}
	}
	return res
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
//

// getbool returns the literal value for a boolean according to the SQL engine.
func (md *metadata) getbool(b bool) string {
	if md.engine == SQLite {
		value := 0
		if b == true {
			value = 1
		}
		return strconv.Itoa(value)
	}

	value := "FALSE"
	if b == true {
		value = "TRUE"
	}
	return value
}

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
