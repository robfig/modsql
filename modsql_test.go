// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package modsql

import (
	"strings"
	"testing"
)

func TestPlaceHolder(t *testing.T) {
	for _, eng := range []Engine{MySQL, Postgres, SQLite} {
		stmtInsert := &Statements{raw: map[int]string{
			0: "INSERT INTO {Q}Foo{Q} (a, b) VALUES({P}, {P})",
		}}
		stmtInsert.setPlaceholder(eng)

		// Check the quote character
		if !strings.Contains(stmtInsert.raw[0], quoteChar[eng]) {
			t.Errorf("%s: expected to get the quote character: %s",
				eng.String(), quoteChar[eng])
		}

		switch eng {
		case MySQL, SQLite:
			if strings.Count(stmtInsert.raw[0], "?") != 2 {
				t.Errorf("expected to get 2 place holders in engine %s: %s",
					eng.String(), stmtInsert.raw[0])
			}
		case Postgres:
			if strings.Count(stmtInsert.raw[0], "$") != 2 {
				t.Errorf("expected to get 2 place holders in engine %s:\n%q",
					eng.String(), stmtInsert.raw[0])
			}
			if !strings.Contains(stmtInsert.raw[0], "$1") &&
				!strings.Contains(stmtInsert.raw[0], "$2") &&
				strings.Contains(stmtInsert.raw[0], "$3") {
				t.Errorf("expected to get place holders correct in engine %s:\n%q",
					eng.String(), stmtInsert.raw[0])
			}
		}
	}
}
