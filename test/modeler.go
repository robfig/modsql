// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	. "github.com/kless/modsql"
	"time"
)

func main() {
	metadata := Metadata("tester", Postgres, MySQL, SQLite)

	Enum("sex", metadata, Int8, 0,
		"female",
		"male",
	)

	types := Table("types", metadata,
		Column("int_", Int).PrimaryKey(),
		Column("int8_", Int8),
		Column("int16_", Int16), //.Validate(),
		Column("int32_", Int32),
		Column("int64_", Int64),

		Column("float32_", Float32),
		Column("float64_", Float64).Index(true),

		Column("string_", String).Unique(),
		Column("binary_", Binary),

		Column("byte_", Byte),
		Column("rune_", Rune).Index(false),
		Column("bool_", Bool),
	)
	types.Unique("float32_", "float64_")
	types.Index(true, "int16_", "int32_")

	def := Table("default_value", metadata,
		Column("id", Int).PrimaryKey(),
		Column("int8_", Int8).Default(int8(55)),
		Column("float32_", Float32).Default(float32(10.2)),

		Column("string_", String),
		Column("binary_", Binary),

		Column("byte_", Byte).Default(byte('b')),
		Column("rune_", Rune).Default('r'),

		Column("bool_", Bool).Default(false),
	)

	times := Table("times", metadata,
		Column("typeId", Int),
		//Column("duration", Duration),
		Column("datetime", DateTime),
	)

	// Insert values

	types.Insert(0, 8, 16, 32, 64, 1.32, 1.64, "one", []byte("12"), 'A', 'Z', true)

	def.InsertTestData(0, 10, 10.10, "foo", []byte{'1', '2'}, 'a', 'z', false)

	times.Insert(0,
		//5*time.Hour+3*time.Minute+12*time.Second,
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	times.Insert(1, time.Time{})

	// == Examples of relationships
	//

	// == Composite foreign keys

	accounts := Table("account", metadata,
		Column("acc_num", Int),
		Column("acc_type", Int),
		Column("acc_descr", String),
	)
	accounts.PrimaryKey("acc_num", "acc_type")

	subAccounts := Table("sub_account", metadata,
		Column("sub_acc", Int).PrimaryKey(),
		Column("ref_num", Int),
		Column("ref_type", Int),
		Column("sub_descr", String),
	)
	subAccounts.Index(false, "ref_num", "ref_type") // MySQL needs individual indexes
	subAccounts.ForeignKey("account",
		ForeignColumn{"ref_num", "acc_num"},
		ForeignColumn{"ref_type", "acc_type"},
	)

	// == One-to-one
	// For related entities which share basic attributes.

	Table("catalog", metadata,
		Column("catalog_id", Int).PrimaryKey(),
		Column("name", String),
		Column("description", String),
		Column("price", Float32),
	)

	Table("magazine", metadata,
		Column("catalog_id", Int).PrimaryKey().ForeignKey("catalog", "catalog_id"),
		Column("page_count", String),
	)

	Table("mp3", metadata,
		Column("catalog_id", Int).PrimaryKey().ForeignKey("catalog", "catalog_id"),
		Column("size", Int),
		Column("length", Float32),
		Column("filename", String),
	)

	// == Many-to-one
	// An item will have many different components, and those components are not
	// of a type that can be shared by other items.

	Table("book", metadata,
		Column("book_id", Int).PrimaryKey(),
		Column("title", String),
		Column("author", String),
	)

	Table("chapter", metadata,
		Column("chapter_id", Int).PrimaryKey(),
		Column("title", String),
		Column("book_fk", Int).ForeignKey("book", "book_id"),
	)

	// == Many-to-many

	// Each user can have several addresses (work, home, grandma's house) and
	// each address can have multiple users.

	Table("user", metadata,
		Column("user_id", Int).PrimaryKey(),
		Column("first_name", String),
		Column("last_name", String),
	)

	Table("address", metadata,
		Column("address_id", Int).PrimaryKey(),
		Column("street", String),
		Column("city", String),
		Column("state", String),
		Column("post_code", String),
	)

	user_addr := Table("user_address", metadata,
		Column("user_id", Int).ForeignKey("user", "user_id"),
		Column("address_id", Int).ForeignKey("address", "address_id"),
	)
	user_addr.PrimaryKey("user_id", "address_id")

	// * * *

	metadata.Create().Write()
}
