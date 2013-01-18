// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build postgresql

package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/bmizerany/pq"
	"github.com/kless/modsql"
	"github.com/kless/modsql/testdata"
)

// To create the database:
//
//   sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//   sudo -u postgres createdb modsql_test --owner USER
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   sudo -u postgres dropdb modsql_test
func TestPostgreSQL(t *testing.T) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbname, host.postgresql))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = modsql.Load(db, "postgresql_init.sql"); err != nil {
		t.Error(err)
	} else {
		if err = modsql.Load(db, "postgresql_test.sql"); err != nil {
			t.Error(err)
		}

		fmt.Println(testdata.ENGINE, "ENG")
	}

	if err = modsql.Load(db, "postgresql_drop.sql"); err != nil {
		t.Error(err)
	}
}
