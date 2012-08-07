// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build sqlite
package modsql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLite(t *testing.T) {
	const dbName = "modsql_test.db"
	defer os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, sqlFilename := getFilenames(SQLite)
	if err = Load(db, sqlFilename); err != nil {
		t.Fatal(err)
	}
}
