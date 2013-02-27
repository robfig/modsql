// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/kless/modsql"
	"github.com/kless/modsql/testdata"
)

// testInsert checks SQL statements generated from Go model.
func testInsert(t *testing.T, db *sql.DB, eng modsql.Engine) {
	testdata.ENGINE = eng
	testdata.Init(db)
	defer testdata.Close()

	err := insertFromTx(db)
	if err != nil {
		testdata.Close()
		t.Fatal(err)
	}

	// Direct data input
	insert := func(model modsql.Modeler) {
		if _, err = model.StmtInsert().Exec(model.Args()...); err != nil {
			t.Error(err)
		}
	}

	insert(&testdata.Types{0, 8, -16, -32, 64, -1.32, -1.64, "a", []byte{1, 2}, 8, 'r', true})
	insert(&testdata.Default_value{0, 8, 1.32, "a", []byte{1, 2}, 8, 'r', false})
	insert(&testdata.Times{0, 7 * time.Hour, time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC)})

	insert(&testdata.Account{11, 22, "a"})
	insert(&testdata.Sub_account{1, 11, 22, "a"})

	insert(&testdata.Catalog{33, "a", "b", 1.32})
	insert(&testdata.Magazine{33, "a"})
	insert(&testdata.Mp3{33, 1, 1.32, "a"})

	insert(&testdata.Book{44, "a", "b"})
	insert(&testdata.Chapter{1, "a", 44})

	insert(&testdata.User{55, "a", "b"})
	insert(&testdata.Address{66, "a", "b", "c", "d"})
	insert(&testdata.User_address{55, 66})
}

// insertFromTx inserts data through a transaction.
func insertFromTx(db *sql.DB) error {
	model := testdata.Catalog{0, "a", "b", 1.32}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Stmt(model.StmtInsert()).Exec(model.Args()...); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
