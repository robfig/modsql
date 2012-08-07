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

// The database was created with:
//
//   sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//   sudo -u postgres createdb modsql_test --owner USER
//
// Note: substitute "USER" by your user name.
func TestPostgreSQL(t *testing.T) {
	const (
		dbname = "modsql_test"
		host   = "/var/run/postgresql"
	)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbname, host))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, sqlFilename := getFilenames(PostgreSQL)
	if err = Load(db, sqlFilename); err != nil {
		t.Fatal(err)
	}
}
