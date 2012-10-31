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
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

var once sync.Once

func init() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")
}

// namesToQuote are names which have to be quoted to be used in SQL statements.
var namesToQuote = [...]string{"user"}

// Queries represent the SQL statements to be used into database functions.
type Queries map[int]string

// Replace replaces the placeholder parameter and adds the quote character
// to tables and columns (if were necessary), according to the SQL engine.
//
// The parameter "from" indicates the SQL engine format for the placeholder
// used in the statement, and "to" is to convert it to another engine format.
//
// TODO: by now, the statement should be in PostgreSQL, so from = PostgreSQL.
func (q Queries) Replace(to, from Engine) {
	rePlaceholder := regexp.MustCompile(`$\d+`) // format in PostgreSQL

	for k, v := range q {
		for _, name := range namesToQuote {
			q[k] = strings.Replace(v, name, quoteChar[to]+name+quoteChar[to], 1)
		}
	}

	if to == MySQL {
		for k, v := range q {
			q[k] = rePlaceholder.ReplaceAllString(v, "?")
		}
	}
}

// * * *

// sqlInt has the integer type for the SQL engine according to the architecture.
var sqlInt = struct {
	MySQLInt   string
	PostgreInt string
}{
	// architecture of 64-bits
	"BIGINT",
	"bigint",
}

// Load loads a database from a file created by ModSQL.
func Load(db *sql.DB, filename string) error {
	once.Do(func() {
		if runtime.GOARCH != "amd64" {
			sqlInt.MySQLInt = "INT"
			sqlInt.PostgreInt = "integer"
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

	for firstLine := ""; ; {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		firstLine += line

		if !strings.HasSuffix(line, ";") {
			continue
		}
		if _, err = db.Exec(firstLine); err != nil {
			return fmt.Errorf("SQL line: %s\n%s", firstLine, err)
		}
		firstLine = ""
	}

	return nil
}

// == Utility
//

// quoteSQL returns the name quoted for SQL.
func quoteSQL(name string) string {
	for _, v := range namesToQuote {
		if v == name {
			return "{{.Q}}" + name + "{{.Q}}"
		}
	}
	return name
}

// quoteSQLField returns field name quoted for SQL.
func quoteSQLField(name string) string {
	for _, v := range namesToQuote {
		if v == name {
			// Add 2 characters by the quotes if are added to the name.
			return "{{.Q}}" + name + "{{.Q}}  "
		}
	}
	return name
}

// validGoName returns a valid name in Go.
func validGoName(name string) string {
	switch name {
	case "type":
		return name + "_"
	}
	return name
}
