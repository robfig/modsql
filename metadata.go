// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"bytes"
	"fmt"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const _FILENAME = "zmodsql"

// metadata defines a collection of table definitions.
type metadata struct {
	useInsert bool
	engines   []sqlEngine
	tables    []*table
	sqlCode   []byte
	goCode    []byte
}

// Metadata returns a new metadata.
func Metadata(eng ...sqlEngine) *metadata {
	for _, v := range eng {
		if err := v.check(); err != nil {
			log.Fatal(err)
		}
	}

	return &metadata{engines: eng}
}

// * * *

// Create generates both SQL statements and Go definitions for all tables.
func (md *metadata) Create() *metadata {
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

	goCode := make([]string, 0)
	sqlCode := make([]string, 0)

	goCode = append(goCode, fmt.Sprintf("%s\npackage %s\n", _HEADER, pkgName))
	sqlCode = append(sqlCode,
		fmt.Sprintf("%s\n%s\nBEGIN;", _CONSTRAINT, _HEADER))

	for _, table := range md.tables {
		// == Get the length of largest field
		fieldMaxLen := 2 // minimum length (id)

		for _, c := range table.columns {
			if len(c.name) > fieldMaxLen {
				fieldMaxLen = len(c.name)
			}
		}
		// ==

		goCode = append(goCode, fmt.Sprintf("\ntype %s struct {\n", table.name))
		sqlCode = append(sqlCode, fmt.Sprintf("\nCREATE TABLE %s (", quote(table.name)))
		columnIndex := make([]string, 0)

		for i, col := range table.columns {
			extra := ""
			field := "\n\t"
			nameQuoted := quote(col.name)

			goCode = append(goCode, fmt.Sprintf("%s %s\n", col.name, col.type_.goString()))

			// == MySQL: Limit the key length in TEXT or BLOB columns
			sqlString := col.type_.tmplAction()

			if col.type_ == String || col.type_ == Binary {
				limit := false

				if col.cons&primaryKey != 0 || col.cons&uniqueCons != 0 {
					limit = true
				}

				if !limit {
					for _, v := range table.uniqueCons {
						if col.name == v {
							limit = true
							break
						}
					}
				}
				if !limit {
					for _, v := range table.pkCons {
						if col.name == v {
							limit = true
							break
						}
					}
				}
				if !limit {
				L:
					for _, fk := range table.fkCons {
						for _, v := range fk.src {
							if col.name == v {
								limit = true
								break L
							}
						}
					}
				}

				if limit {
					sqlString = "VARCHAR(32)"
				}
			}
			// ==

			sqlCode = append(sqlCode, fmt.Sprintf("%s %s%s",
				field+nameQuoted,
				sqlAlign(fieldMaxLen, len(nameQuoted)),
				sqlString,
			))

			if col.cons&primaryKey != 0 {
				extra += " PRIMARY KEY"
			}
			if col.cons&uniqueCons != 0 {
				extra += " UNIQUE"
			}
			if col.cons&foreignKey != 0 {
				extra += fmt.Sprintf(" REFERENCES %s(%s)", col.fkTable, col.fkColumn)
			}

			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch t := col.defaultValue.(type) {
				case bool:
					extra += boolAction(t)
				//case string: extra += fmt.Sprintf("'%s'", t)
				case byte:
					extra += fmt.Sprintf("'%s'", string(t))
				case rune:
					extra += fmt.Sprintf("'%s'", string(t))
				default:
					extra += fmt.Sprintf("%v", t)
				}
			}
			if col.index != 0 {
				unique := ""
				if col.index == uniqIndex {
					unique = "UNIQUE "
				}
				columnIndex = append(columnIndex,
					fmt.Sprintf("CREATE %sINDEX idx_%s_%s ON %s (%s);\n",
						unique, table.name, col.name, table.name, col.name))
			}

			sqlCode = append(sqlCode, extra)

			// The last column
			if i+1 == len(table.columns) {
				var cons []string

				if len(table.uniqueCons) != 0 {
					cons = append(cons, fmt.Sprintf("UNIQUE (%s)",
						strings.Join(table.uniqueCons, ", ")))
				}
				if len(table.pkCons) != 0 {
					cons = append(cons, fmt.Sprintf("PRIMARY KEY (%s)",
						strings.Join(table.pkCons, ", ")))
				}
				for _, fk := range table.fkCons {
					cons = append(cons, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)",
						strings.Join(fk.src, ", "), fk.table,
						strings.Join(fk.dst, ", ")))
				}

				if len(cons) != 0 {
					sqlCode = append(sqlCode, ",\n\n\t"+strings.Join(cons, ",\n\t"))
				}
				sqlCode = append(sqlCode, "\n);\n")
				goCode = append(goCode, "}\n")

				// Indexes
				for i, v := range table.index {
					name := fmt.Sprintf("_m%d", i+1)

					unique := ""
					if v.isUnique {
						unique = "UNIQUE "
					}

					columnIndex = append(columnIndex,
						fmt.Sprintf("CREATE %sINDEX idx_%s_%s ON %s (%s);\n",
							unique, table.name, name, table.name,
							strings.Join(v.index, ", ")))
				}
				if len(columnIndex) != 0 {
					sqlCode = append(sqlCode, columnIndex...)
				}

			} else {
				sqlCode = append(sqlCode, ",")
			}
		}
	}

	// == Insert
	if md.useInsert {
		md.insert(&sqlCode)
	}

	md.sqlCode = []byte(strings.Join(sqlCode, ""))
	md.goCode = []byte(strings.Join(goCode, ""))
	return md
}

// PrintGo prints the Go model.
func (md *metadata) PrintGo() *metadata {
	md.format(os.Stdout)
	return md
}

// PrintSQL prints the SQL statements.
func (md *metadata) PrintSQL() *metadata {
	tmpl, err := template.New("").Parse(string(md.sqlCode))
	if err != nil {
		log.Fatal(err)
	}

	for _, eng := range md.engines {
		if err = tmpl.Execute(os.Stdout, getSQLAction(eng)); err != nil {
			log.Fatal(err)
		}
	}
	return md
}

// Write writes both SQL statements and Go model.
func (md *metadata) Write() {
	if len(md.sqlCode) == 0 {
		log.Fatalf("no data created; use Create()")
	}

	tmpl, err := template.New("").Parse(string(md.sqlCode))
	if err != nil {
		log.Fatal(err)
	}

	for _, eng := range md.engines {
		buf := new(bytes.Buffer)
		if err = tmpl.Execute(buf, getSQLAction(eng)); err != nil {
			log.Fatal(err)
		}

		if err = ioutil.WriteFile(md.sqlFile(eng), buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(_FILENAME+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	md.format(file)
}

// * * *

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
	log.Fatalf("format Go code: %s", err)
}

var replTime = strings.NewReplacer("h", ":", "m", ":", "s", "")

// formatValues converts the values to slice of strings.
func (md *metadata) formatValues(v []interface{}) []string {
	res := make([]string, 0)

	for _, val := range v {
		switch t := val.(type) {
		case bool:
			res = append(res, boolAction(t))
		case int:
			res = append(res, strconv.Itoa(t))
		case float32:
			res = append(res, strconv.FormatFloat(float64(t), 'g', -1, 32))
		case float64:
			res = append(res, strconv.FormatFloat(t, 'g', -1, 64))
		case string, []byte:
			res = append(res, fmt.Sprintf("'%s'", t))
		case time.Duration:
			res = append(res, fmt.Sprintf("'%s'", replTime.Replace(t.String())))
		case time.Time:
			res = append(res, fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05")))
		}
	}
	return res
}

// insert generates SQL statements for insert values; they are finally added to
// the slice main.
func (md *metadata) insert(main *[]string) {
	var data [][]interface{}
	insert := make([]string, 0)

	for _, table := range md.tables {
		tableName := table.name

		data = table.data

		if len(data) != 0 {
			var columns []string

			for _, col := range table.columns {
				columns = append(columns, col.name)
			}
			for _, v := range data {
				insert = append(insert, fmt.Sprintf("\nINSERT INTO %s (%s) VALUES(%s);",
					tableName,
					strings.Join(columns, ", "),
					strings.Join(md.formatValues(v), ", ")))
			}
		}
	}
	insert = append(insert, "\nCOMMIT;\n")

	*main = append(*main, insert...)
}

// sqlFile returns the filename for the SQL statements.
func (md *metadata) sqlFile(eng sqlEngine) string {
	return fmt.Sprintf("%s_%s.sql", _FILENAME, eng.shortString())
}
