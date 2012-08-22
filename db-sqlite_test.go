// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

// +build sqlite

package modsql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLite(t *testing.T) {
	filename := dbname + ".db"
	defer os.Remove(filename)

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, getSQLfile(SQLite)); err != nil {
		t.Fatal(err)
	}
}
