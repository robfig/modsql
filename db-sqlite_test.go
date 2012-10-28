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
	filename := dbname + ".db"
	defer os.Remove(filename)

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, "zsqlite_init.sql"); err != nil {
		t.Error(err)
	} else if err = Load(db, "zsqlite_test.sql"); err != nil {
		t.Error(err)
	}
	if err = Load(db, "zsqlite_drop.sql"); err != nil {
		t.Error(err)
	}
}
