// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql

package main

import (
	"database/sql"
	"fmt"
	"testing"

	//_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/serbaut/go-mysql"
	"github.com/kless/modsql"
)

// To create the database:
//
//   mysql -p
//   mysql> create database modsql_test;
//   mysql> GRANT ALL PRIVILEGES ON modsql_test.* to USER@localhost;
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   mysql> drop database modsql_test;
func TestMySQL(t *testing.T) {
	// Format used in "github.com/Go-SQL-Driver/MySQL"
	//db, err := sql.Open("mysql", fmt.Sprintf("%s@unix(%s)/%s",
	//username, host.mysql, dbname))
	db, err := sql.Open("mysql", fmt.Sprintf("mysql://%s@(unix)/%s?socket=%s",
		username, dbname, host.mysql))
	if err != nil {
		t.Fatal(err)
	}

	if err = modsql.Load(db, "mysql_init.sql"); err != nil {
		t.Error(err)
	} else {
		if err = modsql.Load(db, "mysql_test.sql"); err != nil {
			t.Error(err)
		}

		testInsert(t, db, modsql.MySQL)

		if err = modsql.Load(db, "mysql_drop.sql"); err != nil {
			t.Error(err)
		}
	}

	db.Close()
}
