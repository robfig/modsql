// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Can not pass arguments from tool "go run", so I've to run "go test" for each
// engine.
// build mysql postgresql sqlite

package modsql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestDB(t *testing.T) {
	// Generate files in directory "testdata"
	err := exec.Command("go", "run", "test/example.go").Run()
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Chdir("test"); err != nil {
		log.Fatal(err)
	}

	args := os.Args
	args[0] = "test"
	args = append(args, "-tags")

	args = append(args, "sqlite")
	out, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	args[len(args)-1] = "mysql"
	out, err = exec.Command("go", args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	args[len(args)-1] = "postgresql"
	out, err = exec.Command("go", args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
