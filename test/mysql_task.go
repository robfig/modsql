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

	"github.com/jingweno/gotask/tasking"
	"github.com/kless/modsql"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/serbaut/go-mysql"
)

// NAME
//   test-mysql - check data generated from ModSQL into a MySQL database
//
// DESCRIPTION
//
//   To create the database:
//
//     mysql -p
//     mysql> create database modsql_test;
//     mysql> GRANT ALL PRIVILEGES ON modsql_test.* to USER@localhost;
//
//   Note: substitute "USER" by your user name.
//
//   To remove it:
//
//     mysql> drop database modsql_test;
func TaskTestMySQL(t *tasking.T) {
	// Format used in "github.com/serbaut/go-mysql"
	//db, err := sql.Open("mysql", fmt.Sprintf("mysql://%s@(unix)/%s?socket=%s",
		//username, dbname, host.mysql))
	db, err := sql.Open("mysql", fmt.Sprintf("%s@unix(%s)/%s?parseTime=true",
		username, host.mysql, dbname))
	if err != nil {
		t.Fatal(err)
	}

	if err = modsql.Load(db, filepath.Join("data", "sql", "mysql_init.sql")); err != nil {
		t.Error(err)
	} else {
		if err = modsql.Load(db, filepath.Join("data", "sql", "mysql_test.sql")); err != nil {
			t.Error(err)
		}

		testInsert(t, db, modsql.MySQL)

		if err = modsql.Load(db, filepath.Join("data", "sql", "mysql_drop.sql")); err != nil {
			t.Error(err)
		}
	}

	db.Close()

	if !t.Failed() {
		t.Log("--- PASS")
	}
}
