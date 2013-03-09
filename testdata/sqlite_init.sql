// +build SQLite
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

CREATE TABLE sex (
	id   INTEGER PRIMARY KEY,
	name TEXT
);

CREATE TABLE types (
	int_     INTEGER PRIMARY KEY,
	int8_    INTEGER,
	int16_   INTEGER,
	int32_   INTEGER,
	int64_   INTEGER,
	float32_ REAL,
	float64_ REAL,
	string_  TEXT UNIQUE,
	binary_  BLOB,
	byte_    INTEGER,
	rune_    INTEGER,
	bool_    BOOL,

	UNIQUE (float32_, float64_)
);
CREATE UNIQUE INDEX idx_types_float64_ ON types (float64_);
CREATE INDEX idx_types_rune_ ON types (rune_);
CREATE UNIQUE INDEX idx_types__m1 ON types (int16_, int32_);

CREATE TABLE default_value (
	id       INTEGER PRIMARY KEY,
	int8_    INTEGER DEFAULT 55,
	float32_ REAL DEFAULT 10.2,
	string_  TEXT,
	binary_  BLOB,
	byte_    INTEGER DEFAULT 98,
	rune_    INTEGER DEFAULT 114,
	bool_    BOOL DEFAULT 0
);

CREATE TABLE times (
	typeId   INTEGER,
	datetime TIMESTAMP
);

CREATE TABLE account (
	acc_num   INTEGER,
	acc_type  INTEGER,
	acc_descr TEXT,

	PRIMARY KEY (acc_num, acc_type)
);

CREATE TABLE sub_account (
	sub_acc   INTEGER PRIMARY KEY,
	ref_num   INTEGER,
	ref_type  INTEGER,
	sub_descr TEXT,

	FOREIGN KEY (ref_num, ref_type) REFERENCES account (acc_num, acc_type)
);
CREATE INDEX idx_sub_account__m1 ON sub_account (ref_num, ref_type);

CREATE TABLE catalog (
	catalog_id  INTEGER PRIMARY KEY,
	name        TEXT,
	description TEXT,
	price       REAL
);

CREATE TABLE magazine (
	catalog_id INTEGER PRIMARY KEY REFERENCES catalog(catalog_id),
	page_count TEXT
);

CREATE TABLE mp3 (
	catalog_id INTEGER PRIMARY KEY REFERENCES catalog(catalog_id),
	size       INTEGER,
	length     REAL,
	filename   TEXT
);

CREATE TABLE book (
	book_id INTEGER PRIMARY KEY,
	title   TEXT,
	author  TEXT
);

CREATE TABLE chapter (
	chapter_id INTEGER PRIMARY KEY,
	title      TEXT,
	book_fk    INTEGER REFERENCES book(book_id)
);

CREATE TABLE "user" (
	user_id    INTEGER PRIMARY KEY,
	first_name TEXT,
	last_name  TEXT
);

CREATE TABLE address (
	address_id INTEGER PRIMARY KEY,
	street     TEXT,
	city       TEXT,
	state      TEXT,
	post_code  TEXT
);

CREATE TABLE user_address (
	user_id    INTEGER REFERENCES "user"(user_id),
	address_id INTEGER REFERENCES address(address_id),

	PRIMARY KEY (user_id, address_id)
);

INSERT INTO sex (id, name)
	VALUES(0, 'female');
INSERT INTO sex (id, name)
	VALUES(1, 'male');

INSERT INTO types (int_, int8_, int16_, int32_, int64_, float32_, float64_, string_, binary_, byte_, rune_, bool_)
	VALUES(0, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 65, 90, 1);

INSERT INTO times (typeId, datetime)
	VALUES(0, '2009-11-10T23:00:00Z');
INSERT INTO times (typeId, datetime)
	VALUES(1, NULL);

