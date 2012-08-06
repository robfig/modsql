// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
//	"database/sql"
//	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
//	"testing"

//	_ "github.com/bmizerany/pq"
)

var username string

func init() {
	err := os.Chdir("testdata")
	if err != nil {
		log.Fatal(err)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username = u.Username

	if files, err := filepath.Glob("*.sql"); err != nil {
		log.Fatal(err)
	} else if len(files) == 0 {
		if err = exec.Command("go", "run", "types.go").Run(); err != nil {
			log.Fatal(err)
		}
	}
}

/*func TestPostgre(t *testing.T) {
	const (
		dbName = "modsql_test"
		host   = "/var/run/postgresql"
	)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbName, host))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, PostgreSQL, "./testdata"); err != nil {
		log.Fatal(err)
	}
}*/

func ExampleSQL() {
	metadata := Metadata(SQLite, Help)

	types := Table("types", metadata,
		Column("id", Int64).PrimaryKey(),
		Column("t_int", Int32),
		Column("t_float", Float64),
		Column("t_string", String),
		//Column("t_binary", Binary),
		Column("t_bool", Bool),
	)

	def := Table("default_value", metadata,
		Column("id", Int64).PrimaryKey(),
		Column("d_int", Int8).Default(int8(55)),
		Column("d_float", Float32).Default(float32(10.2)),
		Column("d_string", String).Default("string"),
		//Column("d_binary", Binary).Default([]byte("123")),
		Column("d_bool", Bool).Default(false),
	)

	// == Insert values
	types.InsertHelp("en", "integer", "float", "text", "boolean")
	types.Insert(1, 10, 1.1, "one", true)
	types.Insert(2, 20, 2.2, "two", false)

	def.InsertHelp("en", "integer", "float", "text", "boolean")
	def.Insert(1, 10, 10.1, "foo", true)
	// ==

	metadata.Create().PrintGo().PrintSQL()

	// Output:
/*
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

package modsql

type types struct {
	id       int64
	t_int    int32
	t_float  float64
	t_string string
	t_bool   bool
}

type default_value struct {
	id       int64
	d_int    int8
	d_float  float32
	d_string string
	d_bool   bool
}
// +build sqlite
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN TRANSACTION;

CREATE TABLE types (
	id       INTEGER PRIMARY KEY,
	t_int    INTEGER,
	t_float  REAL,
	t_string TEXT,
	t_bool   BOOL
);

CREATE TABLE _types (
	id       TEXT PRIMARY KEY,
	t_int    TEXT,
	t_float  TEXT,
	t_string TEXT,
	t_bool   TEXT
);

CREATE TABLE default_value (
	id       INTEGER PRIMARY KEY,
	d_int    INTEGER DEFAULT 55,
	d_float  REAL DEFAULT 10.2,
	d_string TEXT DEFAULT 'string',
	d_bool   BOOL DEFAULT 0
);

CREATE TABLE _default_value (
	id       TEXT PRIMARY KEY,
	d_int    TEXT,
	d_float  TEXT,
	d_string TEXT,
	d_bool   TEXT
);

COMMIT;
BEGIN TRANSACTION;

INSERT INTO "_types" (id, t_int, t_float, t_string, t_bool) VALUES('en', 'integer', 'float', 'text', 'boolean');

INSERT INTO "_default_value" (id, d_int, d_float, d_string, d_bool) VALUES('en', 'integer', 'float', 'text', 'boolean');

COMMIT;
BEGIN TRANSACTION;

INSERT INTO "types" (id, t_int, t_float, t_string, t_bool) VALUES(1, 10, 1.1, 'one', 1);
INSERT INTO "types" (id, t_int, t_float, t_string, t_bool) VALUES(2, 20, 2.2, 'two', 0);

INSERT INTO "default_value" (id, d_int, d_float, d_string, d_bool) VALUES(1, 10, 10.1, 'foo', 1);

COMMIT;
*/
}
