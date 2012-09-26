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

var once sync.Once

func init() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")
}

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
