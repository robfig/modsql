// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Can not pass arguments from tool "go run", so I've to run "go test" for each
// engine.
// build mysql postgres sqlite

package modsql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDatabase(t *testing.T) {
	// == Generate files in directory "testdata"
	newTestdata := false
	src := filepath.Join("test", "modeler.go")

	infoSrc, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}
	infoDst, err := os.Stat(filepath.Join("testdata", "model.go"))
	if err != nil {
		newTestdata = true
	}

	if newTestdata || infoSrc.ModTime().After(infoDst.ModTime()) {
		if err = exec.Command("go", "run", src).Run(); err != nil {
			log.Fatal(err)
		}
	}
	//==

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

	args[len(args)-1] = "postgres"
	out, err = exec.Command("go", args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
