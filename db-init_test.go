// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql postgresql sqlite

package modsql

import (
	"log"
	"os"
	"os/exec"
	"os/user"
)

// For access to databases
var (
	host     string
	username string
	dbname   = "modsql_test"
)

func init() {
	err := os.Chdir("testdata")
	if err != nil {
		log.Fatal(err)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username = u.Username

	if err = exec.Command("go", "run", "example.go").Run(); err != nil {
		log.Fatal(err)
	}
}
