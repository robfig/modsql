// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

// namesToQuote are names which have to be quoted to be used in SQL statements
// (tables and columns).
var namesToQuote = [...]string{"user"}

// sqlInt has the integer type for the SQL engine according to the architecture.
// The values could be changed in function Load according to the architecture.
var sqlInt = struct {
	MySQLInt    string
	PostgresInt string
}{
	// architecture of 64-bits
	"BIGINT",
	"bigint",
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("FAIL: ")
}

// Modeler is the interface that wraps the basic Args and StmtInsert methods
// generated in the file "sqlmodel.go".
//
// Args returns the data. It is to be used in prepared statements.
//
// StmtInsert returns the prepared statement to insert data into a later execution.
type Modeler interface {
	Args() []interface{}
	StmtInsert() *sql.Stmt
}

// SQLReplacer replaces "{P}" with the placeholder parameter and "{Q} with
// the quote character, according to the SQL engine.
func SQLReplacer(eng Engine, src string) string {
	switch eng {
	case MySQL, SQLite:
		src = strings.Replace(src, "{P}", "?", -1)

		if strings.Contains(src, "{Q}") {
			return strings.Replace(src, "{Q}", quoteChar[eng], -1)
		}

	case Postgres:
		for nParam := 1; strings.Contains(src, "{P}"); nParam++ {
			src = strings.Replace(src, "{P}", fmt.Sprintf("$%d", nParam), 1)
		}

		if strings.Contains(src, "{Q}") {
			return strings.Replace(src, "{Q}", quoteChar[eng], -1)
		}
	default:
		panic("engine not supported: " + eng.String())
	}
	return src
}

// Statements represents multiple SQL statements prepared to be used with
// different place holders.
type Statements struct {
	raw  map[int]string
	Stmt map[int]*sql.Stmt // to generate from raw
}

// NewStatements returns a set of multiple statements.
// The string to indicate the place holder in raw statements has to be "{P}",
// and the quote character has to be "{Q}".
func NewStatements(raw map[int]string) *Statements {
	return &Statements{
		raw,
		make(map[int]*sql.Stmt, len(raw)),
	}
}

// Prepare creates the prepared statements.
func (m *Statements) Prepare(db *sql.DB, eng Engine) {
	for k, v := range m.raw {
		m.raw[k] = SQLReplacer(eng, v)
	}

	for k, v := range m.raw {
		stmt, err := db.Prepare(v)
		if err != nil {
			log.Fatal(err)
		}
		m.Stmt[k] = stmt
	}
}

// Close closes all prepared statements.
// Returns the first error, if any.
func (m *Statements) Close() error {
	var err, errExit error

	for _, v := range m.Stmt {
		if err = v.Close(); err != nil && errExit == nil {
			errExit = err
		}
	}
	return errExit
}

// listStatements represents a list of Statements to be prepared and closed
// all together.
var listStatements []*Statements

// InitStatements prepares all statements in "listStatements".
// It hast to be called before of insert data.
func InitStatements(db *sql.DB, eng Engine, stmts ...*Statements) {
	listStatements = make([]*Statements, len(stmts))
	for i, v := range stmts {
		v.Prepare(db, eng)
		listStatements[i] = v
	}
}

// CloseStatements closes all statements in "listStatements".
// Returns the first error, if any.
func CloseStatements() error {
	var err, errExit error

	for _, v := range listStatements {
		if err = v.Close(); err != nil && errExit == nil {
			errExit = err
		}
	}
	return errExit
}

// * * *

var once sync.Once

// Load loads a database from a file created by ModSQL.
func Load(db *sql.DB, filename string) error {
	once.Do(func() {
		if runtime.GOARCH != "amd64" {
			sqlInt.MySQLInt = "INT"
			sqlInt.PostgresInt = "integer"
		}
	})

	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, sqlInt); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Handle multiple lines
	for fullLine := ""; ; {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		fullLine += line

		if !strings.HasSuffix(line, ";") { // Multiple line
			continue
		}

		if _, err = tx.Exec(fullLine); err != nil {
			return fmt.Errorf("SQL line: %s\n%s", fullLine, err)
		}
		fullLine = ""
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// == Utility
//

// quoteStatementSQL returns the name quoted for SQL, to use into a statement.
func quoteStatementSQL(name string) string {
	for _, v := range namesToQuote {
		if v == name {
			return "{Q}" + name + "{Q}"
		}
	}
	return name
}

// quoteSQL returns the name quoted for SQL, to use into a template.
func quoteSQL(name string) string {
	for _, v := range namesToQuote {
		if v == name {
			return "{{.Q}}" + name + "{{.Q}}"
		}
	}
	return name
}

// quoteFieldSQL returns field name quoted for SQL, to use into a template.
func quoteFieldSQL(name string) string {
	for _, v := range namesToQuote {
		if v == name {
			// Add 2 characters by the quotes if are added to the name.
			return "{{.Q}}" + name + "{{.Q}}  "
		}
	}
	return name
}
