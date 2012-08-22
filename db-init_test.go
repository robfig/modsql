// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

// +build mysql postgresql sqlite

package modsql

import (
	"log"
	"os"
	"os/exec"
	"os/user"
)

// For access to databases
var (
	host     string
	username string
	dbname   = "modsql_test"
)

func init() {
	err := os.Chdir("testdata")
	if err != nil {
		log.Fatal(err)
	}
	filenameBase = getFileBase() // update by changing of directory

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username = u.Username

	if err = exec.Command("go", "run", "example.go").Run(); err != nil {
		log.Fatal(err)
	}
}
