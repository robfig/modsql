// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gotask

package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jingweno/gotask/tasking"
	"github.com/kless/modsql"
	_ "github.com/mattn/go-sqlite3"
)

// NAME
//   test-sqlite - check data generated from ModSQL into a SQLite database
func TaskTestSQLite(t *tasking.T) {
	filename := dbname + ".db"
	defer func() {
		if err := os.Remove(filename); err != nil {
			t.Error(err)
		}
	}()

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		t.Fatal(err)
	}

	if err = modsql.Load(db, filepath.Join("data", "sql", "sqlite_init.sql")); err != nil {
		t.Error(err)
	} else {
		if err = modsql.Load(db, filepath.Join("data", "sql", "sqlite_test.sql")); err != nil {
			t.Error(err)
		}

		testInsert(t, db, modsql.SQLite)

		if err = modsql.Load(db, filepath.Join("data", "sql", "sqlite_drop.sql")); err != nil {
			t.Error(err)
		}
	}

	db.Close()

	if !t.Failed() {
		t.Log("--- PASS")
	}
}
