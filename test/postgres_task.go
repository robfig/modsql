// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gotask

package main

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/bmizerany/pq"
	"github.com/jingweno/gotask/tasking"
	"github.com/kless/modsql"
)

// NAME
//   check data generated from ModSQL into a Postgre database
//
// DESCRIPTION
//
//   To create the database:
//
//     sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//     sudo -u postgres createdb modsql_test --owner USER
//
//   Note: substitute "USER" by your user name.
//
//   To remove it:
//
//     sudo -u postgres dropdb modsql_test
func TaskTestPostgres(t *tasking.T) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbname, host.postgres))
	if err != nil {
		t.Fatal(err)
	}

	if err = modsql.Load(db, filepath.Join("data", "sql", "postgres_init.sql")); err != nil {
		t.Error(err)
	} else {
		if err = modsql.Load(db, filepath.Join("data", "sql", "postgres_test.sql")); err != nil {
			t.Error(err)
		}

		testInsert(t, db, modsql.Postgres)

		if err = modsql.Load(db, filepath.Join("data", "sql", "postgres_drop.sql")); err != nil {
			t.Error(err)
		}
	}

	db.Close()

	if !t.Failed() {
		t.Log("--- PASS")
	}
}
