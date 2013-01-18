// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql postgresql sqlite

package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// For access to databases

var (
	dbname   = "modsql_test"
	username string
)

var host = struct {
	mysql      string
	postgresql string
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
