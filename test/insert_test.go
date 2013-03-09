// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kless/modsql"
	"github.com/kless/modsql/testdata"
)

// testInsert checks SQL statements generated from Go model.
func testInsert(t *testing.T, db *sql.DB, eng modsql.Engine) {
	modsql.InitStatements(db, eng, testdata.Insert)
	defer modsql.CloseStatements()

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
			in := fmt.Sprintf("%v", input)
			out := fmt.Sprintf("%v", output)

			// The nanoseconds are different in Postgres because it returns fewer digits.
			if eng == modsql.Postgres && strings.Contains(out, "UTC") {
				in = strings.SplitN(out, ".", 2)[0]
				out = strings.SplitN(out, ".", 2)[0]
			}
			if in != out {
				t.Errorf("got different data\ninput:  %v\noutput: %v\n", in, out)
			}
		}
	}

	// Transaction

	inputTx := &testdata.Catalog{0, "a", "b", 1.32}

	err := insertFromTx(db, inputTx)
	if err != nil {
		modsql.CloseStatements()
		t.Fatal(err)
	}

	scan("SELECT * FROM catalog WHERE catalog_id = 0", inputTx, &testdata.Catalog{})

	// Check data input from SQL files

	inputTypes := &testdata.Types{0, 8, 16, 32, 64, 1.32, 1.64, "one", []byte("12"), 'A', 'Z', true}
	scan("SELECT * FROM types WHERE int_ = 0", inputTypes, &testdata.Types{})

	inputDef := &testdata.Default_value{0, 10, 10.10, "foo", []byte{'1', '2'}, 'a', 'z', false}
	scan("SELECT * FROM default_value WHERE Id = 0", inputDef, &testdata.Default_value{})

	inputTimes0 := &testdata.Times{0, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)}
	scan("SELECT * FROM times WHERE typeId = 0", inputTimes0, &testdata.Times{})

	inputTimes1 := &testdata.Times{1, time.Time{}}
	scan("SELECT * FROM times WHERE typeId = 1", inputTimes1, &testdata.Times{})

	// Direct insert

	input0 := &testdata.Types{1, 8, -16, -32, 64, -1.32, -1.64, "a", []byte{1, 2}, 8, 'r', true}
	insert(input0)
	scan("SELECT * FROM types WHERE int_ = 1", input0, &testdata.Types{})

	input1 := &testdata.Default_value{1, 8, 1.32, "a", []byte{1, 2}, 8, 'r', false}
	insert(input1)
	scan("SELECT * FROM default_value WHERE id = 1", input1, &testdata.Default_value{})

	input2 := &testdata.Times{2, time.Now().UTC()}
	insert(input2)
	scan("SELECT * FROM times WHERE typeId = 2", input2, &testdata.Times{})

	input3 := &testdata.Account{11, 22, "a"}
	insert(input3)
	scan("SELECT * FROM account WHERE acc_num = 11", input3, &testdata.Account{})

	input4 := &testdata.Sub_account{1, 11, 22, "a"}
	insert(input4)
	scan("SELECT * FROM sub_account WHERE sub_acc = 1", input4, &testdata.Sub_account{})

	input5 := &testdata.Catalog{33, "a", "b", 1.32}
	insert(input5)
	scan("SELECT * FROM catalog WHERE catalog_id = 33", input5, &testdata.Catalog{})

	input6 := &testdata.Magazine{33, "a"}
	insert(input6)
	scan("SELECT * FROM magazine WHERE catalog_id = 33", input6, &testdata.Magazine{})

	input7 := &testdata.Mp3{33, 1, 1.32, "a"}
	insert(input7)
	scan("SELECT * FROM mp3 WHERE catalog_id = 33", input7, &testdata.Mp3{})

	input8 := &testdata.Book{44, "a", "b"}
	insert(input8)
	scan("SELECT * FROM book WHERE book_id = 44", input8, &testdata.Book{})

	input9 := &testdata.Chapter{1, "a", 44}
	insert(input9)
	scan("SELECT * FROM chapter WHERE chapter_id = 1", input9, &testdata.Chapter{})

	input10 := &testdata.User{55, "a", "b"}
	insert(input10)
	scan("SELECT * FROM {Q}user{Q} WHERE user_id = 55", input10, &testdata.User{})

	input11 := &testdata.Address{66, "a", "b", "c", "d"}
	insert(input11)
	scan("SELECT * FROM address WHERE address_id = 66", input11, &testdata.Address{})

	input12 := &testdata.User_address{55, 66}
	insert(input12)
	scan("SELECT * FROM user_address WHERE user_id = 55", input12, &testdata.User_address{})
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
