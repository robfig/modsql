// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"database/sql"
	"testing"

	"github.com/kless/modsql"
	"github.com/kless/modsql/testdata"
)

// testInsert checks SQL statements generated from Go model.
func testInsert(t *testing.T, db *sql.DB, eng modsql.Engine) {
	testdata.Insert.Prepare(db, eng)
	defer testdata.Insert.Close()

	err := insertFromTx(db)
	if err != nil {
		testdata.Insert.Close()
		t.Fatal(err)
	}

	// Direct data input
	insert := func(model testdata.Modeler) {
		args, err := model.Args()
		if err != nil {
			testdata.Insert.Close()
			t.Fatal(err)
		}
		if _, err = model.StmtInsert().Exec(args...); err != nil {
			t.Error(err)
		}
	}

	insert(&testdata.Catalog{2, "book", "book", 20})
}

// insertFromTx inserts data through a transaction.
func insertFromTx(db *sql.DB) error {
	model := testdata.Catalog{1, "zine", "electronic magazine", 10}

	args, err := model.Args()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Stmt(model.StmtInsert()).Exec(args...); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
