// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gotask

package main

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/jingweno/gotask/tasking"
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
	var t *tasking.T

	u, err := user.Current()
	if err != nil {
		t.Error(err) // Fatal
	}
	username = u.Username

	if err = os.Chdir(filepath.Join("data", "sql")); err != nil {
		t.Error(err) // Fatal
	}
}
