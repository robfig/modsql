// +build SQLite
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE sex (
	id   INTEGER PRIMARY KEY,
	name TEXT
);

CREATE TABLE types (
	t_int     INTEGER PRIMARY KEY,
	t_int8    INTEGER,
	t_int16   INTEGER,
	t_int32   INTEGER,
	t_int64   INTEGER,
	t_float32 REAL,
	t_float64 REAL,
	t_string  TEXT UNIQUE,
	t_binary  BLOB,
	t_byte    TEXT,
	t_rune    TEXT,
	t_bool    BOOL,

	UNIQUE (t_float32, t_float64)
);
CREATE UNIQUE INDEX idx_types_t_float64 ON types (t_float64);
CREATE INDEX idx_types_t_rune ON types (t_rune);
CREATE UNIQUE INDEX idx_types__m1 ON types (t_int16, t_int32);

CREATE TABLE default_value (
	id        INTEGER PRIMARY KEY,
	d_int8    INTEGER DEFAULT 55,
	d_float32 REAL DEFAULT 10.2,
	d_string  TEXT,
	d_binary  BLOB,
	d_byte    TEXT DEFAULT 'b',
	d_rune    TEXT DEFAULT 'r',
	d_bool    BOOL DEFAULT 0
);

CREATE TABLE times (
	typeId     INTEGER,
	t_duration INTEGER,
	t_datetime TEXT
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

INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool)
	VALUES(1, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 'A', 'Z', 1);

INSERT INTO times (typeId, t_duration, t_datetime)
	VALUES(1, '5:3:12', '2009-11-10 23:00:00');
INSERT INTO times (typeId, t_duration, t_datetime)
	VALUES(2, NULL, NULL);
COMMIT;
