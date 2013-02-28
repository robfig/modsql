// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"database/sql"
	"fmt"
	"testing"
	_ "time"

	"github.com/kless/modsql"
	"github.com/kless/modsql/testdata"
)

// testInsert checks SQL statements generated from Go model.
func testInsert(t *testing.T, db *sql.DB, eng modsql.Engine) {
	testdata.ENGINE = eng
	testdata.Init(db)
	defer testdata.Close()

	// insert inserts data without transaction
	insert := func(model modsql.Modeler) {
		if _, err := model.StmtInsert().Exec(model.Args()...); err != nil {
			t.Error(err)
		}
	}

	// scan checks that output data is the same than input data.
	scan := func(query string, input, output modsql.Modeler) {
		rows := db.QueryRow(modsql.SQLReplacer(eng, query))

		if err := rows.Scan(output.Args()...); err != nil {
			t.Error(err)
		} else {
			if fmt.Sprintf("%v", input) != fmt.Sprintf("%v", output) {
				t.Errorf("got different data\ninput:  %v\noutput: %v\n", input, output)
			}
		}
	}

	// Transaction

	inputTx := &testdata.Catalog{0, "a", "b", 1.32}

	err := insertFromTx(db, inputTx)
	if err != nil {
		testdata.Close()
		t.Fatal(err)
	}

	scan("SELECT * FROM catalog WHERE catalog_id = 0", inputTx, &testdata.Catalog{})

	// Direct insert

	input1 := &testdata.Types{0, 8, -16, -32, 64, -1.32, -1.64, "a", []byte{1, 2}, 8, 'r', true}
	insert(input1)
	scan("SELECT * FROM types WHERE t_int = 0", input1, &testdata.Types{})

	input2 := &testdata.Default_value{0, 8, 1.32, "a", []byte{1, 2}, 8, 'r', false}
	insert(input2)
	scan("SELECT * FROM default_value WHERE id = 0", input2, &testdata.Default_value{})

	/*input3 := &testdata.Times{0, 7 * time.Hour, time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC)}
	insert(input3)
	scan("SELECT * FROM times WHERE typeId = 0", input3, &testdata.Times{})*/

	input4 := &testdata.Account{11, 22, "a"}
	insert(input4)
	scan("SELECT * FROM account WHERE acc_num = 11", input4, &testdata.Account{})

	input5 := &testdata.Sub_account{1, 11, 22, "a"}
	insert(input5)
	scan("SELECT * FROM sub_account WHERE sub_acc = 1", input5, &testdata.Sub_account{})

	input6 := &testdata.Catalog{33, "a", "b", 1.32}
	insert(input6)
	scan("SELECT * FROM catalog WHERE catalog_id = 33", input6, &testdata.Catalog{})

	input7 := &testdata.Magazine{33, "a"}
	insert(input7)
	scan("SELECT * FROM magazine WHERE catalog_id = 33", input7, &testdata.Magazine{})

	input8 := &testdata.Mp3{33, 1, 1.32, "a"}
	insert(input8)
	scan("SELECT * FROM mp3 WHERE catalog_id = 33", input8, &testdata.Mp3{})

	input9 := &testdata.Book{44, "a", "b"}
	insert(input9)
	scan("SELECT * FROM book WHERE book_id = 44", input9, &testdata.Book{})

	input10 := &testdata.Chapter{1, "a", 44}
	insert(input10)
	scan("SELECT * FROM chapter WHERE chapter_id = 1", input10, &testdata.Chapter{})

	input11 := &testdata.User{55, "a", "b"}
	insert(input11)
	scan("SELECT * FROM {Q}user{Q} WHERE user_id = 55", input11, &testdata.User{})

	input12 := &testdata.Address{66, "a", "b", "c", "d"}
	insert(input12)
	scan("SELECT * FROM address WHERE address_id = 66", input12, &testdata.Address{})

	input13 := &testdata.User_address{55, 66}
	insert(input13)
	scan("SELECT * FROM user_address WHERE user_id = 55", input13, &testdata.User_address{})
}

// insertFromTx inserts data through a transaction.
func insertFromTx(db *sql.DB, model modsql.Modeler) error {
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
