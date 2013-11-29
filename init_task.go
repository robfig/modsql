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
//   generate files for 'model_task.go'
func TaskInit(t *tasking.T) {
	newTestdata := false
	src := "model_task.go"

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatal(err)
	}
	dstInfo, err := os.Stat(filepath.Join("model", "sqlmodel.go"))
	if err != nil {
		newTestdata = true
	}

	if err := os.Chdir("test"); err != nil {
		t.Fatal(err)
	}
	if newTestdata || srcInfo.ModTime().After(dstInfo.ModTime()) {
		taskBuildModel(t)
	}
}
