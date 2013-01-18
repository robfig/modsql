// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build mysql postgresql sqlite

package modsql

import (
	"log"
	"os/exec"
	"testing"
)

func TestDB(t *testing.T) {
	// Generate files in directory "testdata"
	if err := exec.Command("go", "run", "test/example.go").Run(); err != nil {
		log.Fatal(err)
	}
}
