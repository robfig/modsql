// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gotask

package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jingweno/gotask/tasking"
	"github.com/kless/modsql"
	"github.com/kless/modsql/test/model"
)

// testInsert checks SQL statements generated from Go model.
func testInsert(t *tasking.T, db *sql.DB, eng modsql.Engine) {
	modsql.InitStatements(db, eng, model.Insert)
	defer func() {
		if err := modsql.CloseStatements(); err != nil {
			t.Error(err)
		}
	}()

	// insert inserts data without transaction
	insert := func(model modsql.Modeler) {
		if _, err := model.StmtInsert().Exec(model.Args()...); err != nil {
			t.Error(err)
		}
	}

	// To remove nanoseconds in timestamps since the drivers return fewer digits.
	nsec := regexp.MustCompilePOSIX(`\.[0-9]+ \+`)

	// scan checks that output data is the same than input data.
	scan := func(query string, input, output modsql.Modeler) {
		rows := db.QueryRow(modsql.SQLReplacer(eng, query))

		if err := rows.Scan(output.Args()...); err != nil {
			t.Errorf("query: %q\n%s", query, err)
		} else {
			in := fmt.Sprintf("%v", input)
			out := fmt.Sprintf("%v", output)

			if strings.Contains(out, "UTC") { // Field DateTime
				in = nsec.ReplaceAllLiteralString(in, " +")
				out = nsec.ReplaceAllLiteralString(out, " +")
			}

			if in != out {
				t.Errorf("got different data\ninput:  %v\noutput: %v\n", in, out)
			}
		}
	}

	// Transaction

	inputTx := &model.Catalog{0, "a", "b", 1.32}

	err := insertFromTx(db, inputTx)
	if err != nil {
		modsql.CloseStatements()
		t.Error(err)
	}

	scan("SELECT * FROM catalog WHERE catalog_id = 0", inputTx, &model.Catalog{})

	// Check data input from SQL files

	inputTypes := &model.Types{0, 8, 16, 32, 64, 1.32, 1.64, "one", []byte("12"), 'A', 'Z', true}
	scan("SELECT * FROM types WHERE int_ = 0", inputTypes, &model.Types{})

	inputDef := &model.Default_value{0, 10, 10.10, "foo", []byte{'1', '2'}, 'a', 'z', false}
	scan("SELECT * FROM default_value WHERE Id = 0", inputDef, &model.Default_value{})

	inputTimes0 := &model.Times{0, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)}
	scan("SELECT * FROM times WHERE typeId = 0", inputTimes0, &model.Times{})
	if inputTimes0.Datetime.IsZero() {
		t.Error("inputTimes0.Datetime: should not be zero:", inputTimes0.Datetime)
	}

	inputTimes1 := &model.Times{1, time.Time{}}
	scan("SELECT * FROM times WHERE typeId = 1", inputTimes1, &model.Times{})
	if !inputTimes1.Datetime.IsZero() {
		t.Error("inputTimes1.Datetime: should be zero:", inputTimes1.Datetime)
	}

	// Direct insert

	input0 := &model.Types{1, 8, -16, -32, 64, -1.32, -1.64, "a", []byte{1, 2}, 8, 'r', true}
	insert(input0)
	scan("SELECT * FROM types WHERE int_ = 1", input0, &model.Types{})

	input1 := &model.Default_value{1, 8, 1.32, "a", []byte{1, 2}, 8, 'r', false}
	insert(input1)
	scan("SELECT * FROM default_value WHERE id = 1", input1, &model.Default_value{})

	input2 := &model.Times{2, time.Now().UTC()}
	insert(input2)
	scan("SELECT * FROM times WHERE typeId = 2", input2, &model.Times{})
	if input2.Datetime.IsZero() {
		t.Error("input2.Datetime: should not be zero:", input2.Datetime)
	}

	input3 := &model.Account{11, 22, "a"}
	insert(input3)
	scan("SELECT * FROM account WHERE acc_num = 11", input3, &model.Account{})

	input4 := &model.Sub_account{1, 11, 22, "a"}
	insert(input4)
	scan("SELECT * FROM sub_account WHERE sub_acc = 1", input4, &model.Sub_account{})

	input5 := &model.Catalog{33, "a", "b", 1.32}
	insert(input5)
	scan("SELECT * FROM catalog WHERE catalog_id = 33", input5, &model.Catalog{})

	input6 := &model.Magazine{33, "a"}
	insert(input6)
	scan("SELECT * FROM magazine WHERE catalog_id = 33", input6, &model.Magazine{})

	input7 := &model.Mp3{33, 1, 1.32, "a"}
	insert(input7)
	scan("SELECT * FROM mp3 WHERE catalog_id = 33", input7, &model.Mp3{})

	input8 := &model.Book{44, "a", "b"}
	insert(input8)
	scan("SELECT * FROM book WHERE book_id = 44", input8, &model.Book{})

	input9 := &model.Chapter{1, "a", 44}
	insert(input9)
	scan("SELECT * FROM chapter WHERE chapter_id = 1", input9, &model.Chapter{})

	input10 := &model.User{55, "a", "b"}
	insert(input10)
	scan("SELECT * FROM {Q}user{Q} WHERE user_id = 55", input10, &model.User{})

	input11 := &model.Address{66, "a", "b", "c", "d"}
	insert(input11)
	scan("SELECT * FROM address WHERE address_id = 66", input11, &model.Address{})

	input12 := &model.User_address{55, 66}
	insert(input12)
	scan("SELECT * FROM user_address WHERE user_id = 55", input12, &model.User_address{})
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
