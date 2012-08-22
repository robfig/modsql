// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

// +build postgresql

package modsql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/bmizerany/pq"
)

// To create the database:
//
//   sudo -u postgres createuser USER --no-superuser --no-createrole --no-createdb
//   sudo -u postgres createdb modsql_test --owner USER
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   sudo -u postgres dropdb modsql_test
func TestPostgreSQL(t *testing.T) {
	host = "/var/run/postgresql"

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable",
		username, dbname, host))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, getSQLfile(PostgreSQL)); err != nil {
		t.Fatal(err)
	}
}
