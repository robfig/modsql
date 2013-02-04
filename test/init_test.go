// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql postgres sqlite

package main

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kless/modsql/testdata"
)

// For access to databases

var (
	dbname   = "modsql_test"
	username string
)

var host = struct {
	mysql    string
	postgres string
}{
	"/var/run/mysqld/mysqld.sock",
	"/var/run/postgresql",
}

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username = u.Username

	if err = os.Chdir(filepath.Join("..", "testdata")); err != nil {
		log.Fatal(err)
	}
}

// * * *

// testFromModel checks SQL statements generated from the Go model.
func testFromModel(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	m := testdata.Catalog{1, "zine", "electronic magazine", 10}
	if _, err = db.Exec(m.Insert()); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
