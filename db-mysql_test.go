// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql
package modsql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "code.google.com/p/go-mysql-driver/mysql"
)

// The database was created with:
//
//   mysql -p
//   mysql> create database modsql_test;
//   mysql> GRANT ALL PRIVILEGES ON modsql_test.* to neo@localhost;
//
// Note: substitute "neo" by your user name.
func TestMySQL(t *testing.T) {
	const (
		dbname = "modsql_test"
		host   = "/var/run/mysqld/mysqld.sock"
	)

	db, err := sql.Open("mysql", fmt.Sprintf("%s@unix(%s)/%s?charset=utf8",
		username, host, dbname))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, sqlFilename := getFilenames(MySQL)
	if err = Load(db, sqlFilename); err != nil {
		t.Fatal(err)
	}
}
