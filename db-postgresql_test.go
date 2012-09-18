// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build postgresql

package modsql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/bmizerany/pq"
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
	host = "/var/run/postgresql"

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbname, host))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, "zmodsql_pg.sql"); err != nil {
		t.Fatal(err)
	}
}
