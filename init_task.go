// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gotask

package modsql

import (
	"os"
	"path/filepath"

	"github.com/jingweno/gotask/tasking"
)

// NAME
//   generate files for 'test/modeler.go'
func TaskInit(t *tasking.T) {
	if err := os.Chdir("test"); err != nil {
		t.Fatal(err)
	}

	newTestdata := false
	src := "modeler.go"

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatal(err)
	}
	dstInfo, err := os.Stat(filepath.Join("tester", "sqlmodel.go"))
	if err != nil {
		newTestdata = true
	}

	if newTestdata || srcInfo.ModTime().After(dstInfo.ModTime()) {
		if err = t.Exec("go", "run", src); err != nil {
			t.Error(err)
		}
	}
}
