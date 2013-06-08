// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

const (
	_CONSTRAINT  = "// +build {{.Engine}}"
	_HEADER      = "// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT\n"
	_HEADER_EDIT = "// MACHINE GENERATED BY ModSQL (github.com/kless/modsql)\n"
)

// metadata defines a collection of table definitions.
type metadata struct {
	useInsert     bool
	useInsertTest bool

	posQueries int

	engines []Engine
	tables  []*table

	goCode    []string
	sqlCreate []string
	sqlDrop   []string
	sqlTest   []string
	sqlInsert []string

	dirPath string
	pkgName string
}

// Metadata returns a new metadata.
// The package name is used to create the directory where to generate all files
// for every engine and for the Go's package.
//
// The new directory is created in the path where it is run.
func Metadata(packageName string, eng ...Engine) *metadata {
	for _, v := range eng {
		if err := v.check(); err != nil {
			log.Fatal(err)
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir = filepath.Join(dir, packageName)
	if err = os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	return &metadata{engines: eng, dirPath: dir, pkgName: packageName}
}

// * * *

// Create generates both SQL statements and Go definitions for all tables.
func (md *metadata) Create() *metadata {
	// Align SQL types adding spaces.
	sqlAlign := func(maxLen, nameLen int) string {
		if maxLen <= nameLen {
			return ""
		}
		return strings.Repeat(" ", maxLen-nameLen)
	}

	md.goCode = append(md.goCode, fmt.Sprintf("%s\npackage %s\n", _HEADER_EDIT, md.pkgName))

	md.goCode = append(md.goCode, "import (\n\"database/sql\"\n")
	md.goCode = append(md.goCode, "") // Could add another import
	md.goCode = append(md.goCode, "\n\"github.com/kless/modsql\"\n)")

	md.goCode = append(md.goCode, "\n// == EDIT\n")
	for i, v := range md.engines {
		if i != 0 {
			md.goCode = append(md.goCode, "//")
		}
		md.goCode = append(md.goCode, "const ENGINE = modsql."+v.String()+"\n")
	}

	md.goCode = append(md.goCode, `
// * * *

var Insert = modsql.NewStatements(map[int]string{
`)

	// To add all queries
	md.posQueries = len(md.goCode)
	md.goCode = append(md.goCode, "")

	md.goCode = append(md.goCode, `,
})

`)

	md.sqlCreate = append(md.sqlCreate,
		fmt.Sprintf("%s\n%s", _CONSTRAINT, _HEADER))

	md.sqlDrop = append(md.sqlDrop, _HEADER)
	md.sqlDrop = append(md.sqlDrop, "{{.MySQLDrop0}}")

	useTime := false
	iTable := 0 // to differenciate from tables for enums

	for _, table := range md.tables {
		// == Get the length of largest field
		fieldMaxLen := 2 // minimum length (id)

		for _, c := range table.columns {
			if len(c.name) > fieldMaxLen {
				fieldMaxLen = len(c.name)
			}
		}
		// ==

		if !table.isEnum {
			md.goCode = append(md.goCode,
				fmt.Sprintf("\ntype %s struct {\n", strings.Title(table.name)))
		} else {
			md.goCode = append(md.goCode, "\n// "+table.name+"\nconst(\n")
		}

		md.sqlCreate = append(md.sqlCreate,
			fmt.Sprintf("\nCREATE TABLE %s (", table.sqlName))
		md.sqlDrop = append(md.sqlDrop,
			fmt.Sprintf("\nDROP TABLE %s{{.PostgresDrop}};", table.sqlName))

		columnIndex := make([]string, 0)
		columnNames := make([]string, 0)
		columnValues := make([]string, 0)

		for iCol, col := range table.columns {
			extra := ""

			//if !useTime && (col.type_ == Duration || col.type_ == DateTime) {
			if !useTime && col.type_ == DateTime {
				useTime = true
			}

			if !table.isEnum {
				type_ := col.type_.goString()

				md.goCode = append(md.goCode,
					fmt.Sprintf("%s %s\n", strings.Title(col.name), type_))
				columnNames = append(columnNames, col.name)
				columnValues = append(columnValues, type_)
			} else if iCol == 0 {
				name := table.name

				// Get the first part of the table name; until '_' or letter is upper
				for iName, letter := range table.name[1:] {
					if unicode.IsUpper(letter) || letter == '_' {
						name = table.name[:iName+1]
						break
					}
				}
				name = strings.ToUpper(name) + "_"

				for iData, vData := range table.data {
					if iData == 0 {
						iota_ := "iota"

						if table.startEnum != 0 {
							iota_ += " + " + strconv.Itoa(table.startEnum)
						}
						md.goCode = append(md.goCode, fmt.Sprintf("%s = %s\n",
							name+strings.ToUpper(vData[1].(string)), iota_))
					} else {
						md.goCode = append(md.goCode, fmt.Sprintf("%s\n",
							name+strings.ToUpper(vData[1].(string))))
					}
				}
			}

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
					sqlString = "{{.StringLimit}}"
				}
			}
			// ==
			nameQuoted := quoteFieldSQL(col.name)

			md.sqlCreate = append(md.sqlCreate, fmt.Sprintf("\n\t%s %s%s",
				nameQuoted, sqlAlign(fieldMaxLen, len(nameQuoted)), sqlString))

			if col.cons&primaryKey != 0 {
				extra += " PRIMARY KEY"
			}
			if col.cons&uniqueCons != 0 {
				extra += " UNIQUE"
			}
			if col.cons&foreignKey != 0 {
				extra += fmt.Sprintf(" REFERENCES %s(%s)", quoteSQL(col.fkTable), col.fkColumn)
			}

			if col.defaultValue != nil {
				extra += " DEFAULT "

				switch t := col.defaultValue.(type) {
				case bool:
					extra += boolAction(t)
				//case string: extra += fmt.Sprintf("'%s'", t)
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
						unique, table.name, col.name, table.sqlName, col.name))
			}

			md.sqlCreate = append(md.sqlCreate, extra)

			// The last column
			if iCol+1 == len(table.columns) {
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
						strings.Join(fk.src, ", "), quoteSQL(fk.table),
						strings.Join(fk.dst, ", ")))
				}

				if len(cons) != 0 {
					md.sqlCreate = append(md.sqlCreate, ",\n\n\t"+strings.Join(cons, ",\n\t"))
				}
				md.sqlCreate = append(md.sqlCreate, "\n);\n")
				if !table.isEnum {
					md.goCode = append(md.goCode, "}\n")

					md.goCode = append(md.goCode,
						md.genInsertForType(iTable, table.name, columnNames, columnValues),
					)
					iTable++
				} else {
					md.goCode = append(md.goCode, ")\n")
				}

				// Indexes
				for i, v := range table.index {
					name := fmt.Sprintf("_m%d", i+1)

					unique := ""
					if v.isUnique {
						unique = "UNIQUE "
					}

					columnIndex = append(columnIndex,
						fmt.Sprintf("CREATE %sINDEX idx_%s_%s ON %s (%s);\n",
							unique, table.name, name, table.sqlName,
							strings.Join(v.index, ", ")))
				}
				if len(columnIndex) != 0 {
					md.sqlCreate = append(md.sqlCreate, columnIndex...)
				}

			} else {
				md.sqlCreate = append(md.sqlCreate, ",")
			}
		}
	}

	if useTime {
		md.goCode[2] = "\"time\"\n"
	}

	md.goCode[md.posQueries] = strings.Join(md.sqlInsert, ",\n")

	// == Insert
	if md.useInsert {
		md.sqlCreate = append(md.sqlCreate, md.genInsert(false)...)
	}
	if md.useInsertTest {
		md.sqlTest = append(md.sqlTest, _HEADER)
		md.sqlTest = append(md.sqlTest, md.genInsert(true)...)
	}
	md.sqlDrop = append(md.sqlDrop, "{{.MySQLDrop1}}")
	md.sqlDrop = append(md.sqlDrop, "\n\n")

	return md
}

// PrintGo prints the Go model.
func (md *metadata) PrintGo() *metadata {
	md.format(os.Stdout)
	return md
}

// PrintSQL prints the SQL statements.
func (md *metadata) PrintSQL() *metadata {
	tmpl, err := template.New("").Parse(strings.Join(md.sqlCreate, ""))
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
	if len(md.sqlCreate) == 0 {
		log.Fatalf("no data created; use Create()")
	}

	tmplCreate, err := template.New("").Parse(strings.Join(md.sqlCreate, ""))
	if err != nil {
		log.Fatal(err)
	}
	tmplDrop, err := template.New("").Parse(strings.Join(md.sqlDrop, ""))
	if err != nil {
		log.Fatal(err)
	}

	var tmplTest *template.Template
	if md.useInsertTest {
		tmplTest, err = template.New("").Parse(strings.Join(md.sqlTest, ""))
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, eng := range md.engines {
		filename := filepath.Join(md.dirPath, strings.ToLower(eng.String()))

		buf := new(bytes.Buffer)
		if err = tmplCreate.Execute(buf, getSQLAction(eng)); err != nil {
			log.Fatal(err)
		}
		if err = ioutil.WriteFile(filename+"_init.sql", buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

		buf = new(bytes.Buffer)
		if err = tmplDrop.Execute(buf, getSQLAction(eng)); err != nil {
			log.Fatal(err)
		}
		if err = ioutil.WriteFile(filename+"_drop.sql", buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

		if md.useInsertTest {
			buf = new(bytes.Buffer)
			if err = tmplTest.Execute(buf, getSQLAction(eng)); err != nil {
				log.Fatal(err)
			}
			if err = ioutil.WriteFile(filename+"_test.sql", buf.Bytes(), 0644); err != nil {
				log.Fatal(err)
			}
		}
	}

	file, err := os.OpenFile(filepath.Join(md.dirPath, "model.go"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	md.format(file)
}

// * * *

// TimeReplacer is a replacer for duration times.
// It is to be called from the Go code generated.
var TimeReplacer = strings.NewReplacer("h", ":", "m", ":", "s", "")

// format formats the Go source code.
func (md *metadata) format(out io.Writer) {
	codeFmt, err := format.Source([]byte(strings.Join(md.goCode, "")))
	if err != nil {
		goto _error
	}
	if _, err = out.Write(codeFmt); err != nil {
		goto _error
	}

	return
_error:
fmt.Println(strings.Join(md.goCode, "")) // TODO: remove
	log.Fatalf("format Go code: %s", err)
}

// genInsert generates SQL statements to insert values.
// If testdata is true, it generates values for test.
func (md *metadata) genInsert(testdata bool) []string {
	var data [][]interface{}
	insert := make([]string, 0)

	for _, table := range md.tables {
		if testdata {
			data = table.testData
		} else {
			data = table.data
		}

		if len(data) != 0 {
			var columns []string

			for _, col := range table.columns {
				columns = append(columns, quoteSQL(col.name))
			}
			for _, v := range data {
				insert = append(insert, fmt.Sprintf("\nINSERT INTO %s (%s)\n\tVALUES(%s);",
					table.sqlName,
					strings.Join(columns, ", "),
					formatSQL(v)))
			}
			insert = append(insert, "\n")
		}
	}
	insert = append(insert, "\n")
	return insert
}

// * * *

// formatSQL converts the values to a string formatted in SQL.
func formatSQL(v []interface{}) string {
	res := make([]string, len(v))

	for i, val := range v {
		switch t := val.(type) {
		case bool:
			res[i] = boolAction(t)

		case int:
			res[i] = strconv.Itoa(t)
		case int8:
			res[i] = strconv.Itoa(int(t))
		case int16:
			res[i] = strconv.Itoa(int(t))
		case int32:
			res[i] = strconv.Itoa(int(t))
		case int64:
			res[i] = strconv.Itoa(int(t))

		case float32:
			res[i] = strconv.FormatFloat(float64(t), 'g', -1, 32)
		case float64:
			res[i] = strconv.FormatFloat(t, 'g', -1, 64)

		case string, []byte:
			res[i] = fmt.Sprintf("'%s'", t)

		//case time.Duration:
			//res[i] = fmt.Sprintf("'%s'", TimeReplacer.Replace(t.String()))
		case time.Time:
			res[i] = fmt.Sprintf("'%s'", t.Format(time.RFC3339Nano))

		case nil:
			res[i] = "NULL"
		}
	}
	return strings.Join(res, ", ")
}

// genInsertForType generate the SQL statement to insert data from a Go type.
func (md *metadata) genInsertForType(idx int, name string, columns, values []string) string {
	args := make([]string, len(columns))

	for i, _ := range values {
		addColumn := true

		/*switch v {
		//case "time.Duration":
			//addColumn = false
			//args[i] = fmt.Sprintf("modsql.TimeReplacer.Replace(t.%s.String())", strings.Title(columns[i]))
		case "time.Time":
			addColumn = false
			args[i] = fmt.Sprintf("t.%s.Format(time.RFC3339Nano)", strings.Title(columns[i]))
		}*/

		if addColumn {
			args[i] = "&t." + strings.Title(columns[i])
		}
	}

	tmplArgs := strings.Repeat("{P}, ", len(args))
	md.sqlInsert = append(md.sqlInsert,
		fmt.Sprintf("%d: \"INSERT INTO %s (%s) VALUES(%s)\"",
			len(md.sqlInsert), quoteStatementSQL(name), strings.Join(columns, ", "),
			tmplArgs[:len(tmplArgs)-2]),
	)

	name = strings.Title(name)

	return fmt.Sprintf(
		"func (t *%s) Args() []interface{} {\n"+
			"return []interface{}{%s}\n"+
		"}\n\n"+

		"func (t *%s) StmtInsert() *sql.Stmt { return Insert.Stmt[%d] }",

		name,
		strings.Join(args, ", "),
		name,
		idx,
	)
}

/*func (t *Types) Validate() error {
	v := validate.NewRFCEmail(validate.TrimSpace)
	t.T_string, err := v.Get(t.T_string)
	if err != nil {
		return err
	}
}*/
