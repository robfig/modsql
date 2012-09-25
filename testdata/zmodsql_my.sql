// +build mysql
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE types (
	t_int     {{.MySQLInt}} PRIMARY KEY,
	t_int8    TINYINT,
	t_int16   SMALLINT,
	t_int32   INT,
	t_int64   BIGINT,
	t_float32 FLOAT,
	t_float64 DOUBLE,
	t_string  VARCHAR(255) UNIQUE,
	t_binary  BLOB,
	t_byte    CHAR(1),
	t_rune    CHAR(4),
	t_bool    BOOL,

	UNIQUE (t_float32, t_float64)
);
CREATE UNIQUE INDEX idx_types_t_float64 ON types (t_float64);
CREATE INDEX idx_types_t_rune ON types (t_rune);
CREATE INDEX idx_types__m1 ON types (t_int8, t_float32);
CREATE UNIQUE INDEX idx_types__m2 ON types (t_int16, t_int32);

CREATE TABLE default_value (
	id        {{.MySQLInt}} PRIMARY KEY,
	d_int8    TINYINT DEFAULT 55,
	d_float32 FLOAT DEFAULT 10.2,
	d_string  TEXT,
	d_binary  BLOB,
	d_byte    CHAR(1) DEFAULT 'b',
	d_rune    CHAR(4) DEFAULT 'r',
	d_bool    BOOL DEFAULT FALSE
);

CREATE TABLE times (
	typeId     {{.MySQLInt}} REFERENCES types(t_int),
	t_duration TIME,
	t_datetime TIMESTAMP,

	PRIMARY KEY (t_duration, t_datetime)
);
CREATE INDEX idx_times_t_datetime ON times (t_datetime);

CREATE TABLE account (
	acc_num   {{.MySQLInt}},
	acc_type  {{.MySQLInt}},
	acc_descr TEXT,

	PRIMARY KEY (acc_num, acc_type)
);

CREATE TABLE sub_account (
	sub_acc   {{.MySQLInt}} PRIMARY KEY,
	ref_num   {{.MySQLInt}},
	ref_type  {{.MySQLInt}},
	sub_descr TEXT,

	FOREIGN KEY (ref_type, ref_num) REFERENCES account (acc_type, acc_num)
);

CREATE TABLE catalog (
	catalog_id  {{.MySQLInt}} PRIMARY KEY,
	name        TEXT,
	description TEXT,
	price       FLOAT
);

CREATE TABLE magazine (
	catalog_id {{.MySQLInt}} PRIMARY KEY REFERENCES catalog(catalog_id),
	page_count TEXT
);

CREATE TABLE mp3 (
	catalog_id {{.MySQLInt}} PRIMARY KEY REFERENCES catalog(catalog_id),
	size       {{.MySQLInt}},
	length     FLOAT,
	filename   TEXT
);

CREATE TABLE book (
	book_id {{.MySQLInt}} PRIMARY KEY,
	title   TEXT,
	author  TEXT
);

CREATE TABLE chapter (
	chapter_id {{.MySQLInt}} PRIMARY KEY,
	title      TEXT,
	book_fk    {{.MySQLInt}} REFERENCES book(book_id)
);

CREATE TABLE person (
	person_id  {{.MySQLInt}} PRIMARY KEY,
	first_name TEXT,
	last_name  TEXT
);

CREATE TABLE address (
	address_id {{.MySQLInt}} PRIMARY KEY,
	street     TEXT,
	city       TEXT,
	state      TEXT,
	post_code  TEXT
);

CREATE TABLE person_address (
	person_id  {{.MySQLInt}} REFERENCES person(person_id),
	address_id {{.MySQLInt}} REFERENCES address(address_id),

	PRIMARY KEY (person_id, address_id)
);

INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool) VALUES(1, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 'A', 'Z', TRUE);
INSERT INTO default_value (id, d_int8, d_float32, d_string, d_binary, d_byte, d_rune, d_bool) VALUES(1, 10, 10.1, 'foo', '12', 'a', 'z', FALSE);
INSERT INTO times (typeId, t_duration, t_datetime) VALUES(1, '5:3:12', '2012-09-25 07:48:17');
COMMIT;
