// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

// +build mysql
package modsql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "code.google.com/p/go-mysql-driver/mysql"
)

// To create the database:
//
//   mysql -p
//   mysql> create database modsql_test;
//   mysql> GRANT ALL PRIVILEGES ON modsql_test.* to USER@localhost;
//
// Note: substitute "USER" by your user name.
//
// To remove it:
//
//   mysql> drop database modsql_test;
func TestMySQL(t *testing.T) {
	host = "/var/run/mysqld/mysqld.sock"

	db, err := sql.Open("mysql", fmt.Sprintf("%s@unix(%s)/%s?charset=utf8",
		username, host, dbname))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Load(db, getSQLfile(MySQL)); err != nil {
		t.Fatal(err)
	}
}
