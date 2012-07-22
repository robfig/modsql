// Copyright 2010  The "ModSQL" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
	"strconv"
	"strings"
)

// mode represents the modes to use in metadata.Mode.
type mode byte

const (
	// If Help is set, it is created tables related to help.
	Help mode = iota + 1
)

// metadata defines a collection of table definitions.
type metadata struct {
	engine        sqlEngine
	mode          mode
	useInsert     bool
	useInsertHelp bool
	tables        []*table
	sqlCode       []byte
	goCode        []byte
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
	sqlCode := make([]string, 0, 0)
	goCode := make([]string, 0, 0)

	pop := func(sl []string) []string {
		_, sl = sl[len(sl)-1], sl[:len(sl)-1]
		return sl
	}

	// Quote special names.
	quote := func(name string) string {
		if name == "user" {
			return `"` + name + `"`
		}
		return name
	}

	// Align SQL types adding spaces.
	sqlAlign := func(maxLen, nameLen int) string {
		if maxLen == nameLen {
			return ""
		}
		return strings.Repeat(" ", maxLen-nameLen)
	}

	// Package name
	pkgName := "main"
	pkg, err := build.ImportDir(".", 0)
	if err == nil {
		pkgName = pkg.Name
	}

	goCode = append(goCode, fmt.Sprintf("%s\npackage %s\n", header, pkgName))
	sqlCode = append(sqlCode, fmt.Sprintf("%s\nBEGIN TRANSACTION;\n", header))

	for _, table := range md.tables {
		sqlLangCode := make([]string, 0, 0)

		// == Get the length of largest field
		fieldMaxLen := 2 // minimum length (id)

		for _, c := range table.columns {
			if len(c.name) > fieldMaxLen {
				fieldMaxLen = len(c.name)
			}
		}
		// ==

		if md.mode == Help {
			sqlLangCode = append(sqlLangCode,
				fmt.Sprintf("\nCREATE TABLE _%s (\n\tid %sTEXT PRIMARY KEY,\n",
					table.name, sqlAlign(fieldMaxLen, 2)))
		}

		goCode = append(goCode, fmt.Sprintf("\ntype %s struct {\n", table.name))
		sqlCode = append(sqlCode, fmt.Sprintf("\nCREATE TABLE %s (", quote(table.name)))

		for i, col := range table.columns {
			extra := ""
			field := "\n\t"
			nameQuoted := quote(col.name)

			goCode = append(goCode, fmt.Sprintf("%s %s\n", col.name, col.type_.goString()))

			sqlCode = append(sqlCode, fmt.Sprintf("%s %s%s",
				field+nameQuoted,
				sqlAlign(fieldMaxLen, len(nameQuoted)),
				col.type_.sqlString(md.engine),
			))

			if col.isPrimaryKey {
				extra += " PRIMARY KEY"
			}
			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch t := col.defaultValue.(type) {
				case string:
					extra += fmt.Sprintf("'%s'", t)
				case bool:
					extra += md.formatBool(t)
				default:
					extra += fmt.Sprintf("%v", t)
				}
			}

			sqlCode = append(sqlCode, extra)

			// Add table for translation of fields comments
			if md.mode == Help && col.name != "id" {
				sqlLangCode = append(sqlLangCode, fmt.Sprintf("\t%s %sTEXT",
					nameQuoted, sqlAlign(fieldMaxLen, len(nameQuoted))),
				)
				sqlLangCode = append(sqlLangCode, ",\n")
			}

			// The last column
			if i+1 == len(table.columns) {
				sqlCode = append(sqlCode, "\n);\n")
				goCode = append(goCode, "}\n")

				if md.mode == Help {
					sqlLangCode = pop(sqlLangCode)
					sqlLangCode = append(sqlLangCode, "\n);\n")
					sqlCode = append(sqlCode, sqlLangCode...)
				}
			} else {
				sqlCode = append(sqlCode, ",")
			}
		}
	}
	sqlCode = append(sqlCode, "\nCOMMIT;\n")

	// == Insert
	if md.useInsertHelp {
		md.insert(&sqlCode, _INSERT_HELP)
	}
	if md.useInsert {
		md.insert(&sqlCode, _INSERT_DATA)
	}

	md.sqlCode = []byte(strings.Join(sqlCode, ""))
	md.goCode = []byte(strings.Join(goCode, ""))
	return md
}

// Print prints both SQL statements and Go model.
func (md *metadata) Print() {
	fmt.Printf("%s\n* * *\n\n", md.sqlCode)
	md.format(os.Stdout)
}

// Write writes both SQL statements and Go model to files using names by default.
func (md *metadata) Write() {
	md.WriteTo(_SQL_FILE, _MODEL_FILE)
}

// WriteTo writes both SQL statements and Go model to given files.
func (md *metadata) WriteTo(sqlFile, goFile string) {
	if len(md.sqlCode) == 0 {
		_log.Fatalf("no tables created; use CreateAll()")
	}

	err := ioutil.WriteFile(sqlFile, md.sqlCode, 0644)
	if err != nil {
		_log.Fatalf("write file: %s", err)
	}

	file, err := os.OpenFile(goFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		_log.Fatalf("open file: %s", err)
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

// format formats the Go source code.
func (md *metadata) format(out io.Writer) {
	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", md.goCode, _PARSER_MODE)
	if err != nil {
		goto _error
	}

	err = (&printer.Config{_PRINTER_MODE, _TAB_WIDTH}).Fprint(out, fset, ast)
	if err != nil {
		goto _error
	}

	return
_error:
	_log.Fatalf("format Go code: %s", err)
}

// insert generates SQL statements to insert values; they are finally added to
// the slice main.
func (md *metadata) insert(main *[]string, value uint) {
	if value != _INSERT_HELP && value != _INSERT_DATA {
		_log.Fatalf("argument \"value\" not valid for \"metadata.insert\": %d", value)
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
					strings.Join(md.formatValues(v), ", ")))
			}
			insert = append(insert, "\n")
		}
	}

	insert = append(insert, "\nCOMMIT;\n")
	*main = append(*main, insert...)
}

// formatValues converts the values to slice of strings.
func (md *metadata) formatValues(v []interface{}) []string {
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
			res = append(res, md.formatBool(t))
		}
	}
	return res
}
