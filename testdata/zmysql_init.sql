// +build MySQL
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE sex (
	id   TINYINT PRIMARY KEY,
	name TEXT
);

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
CREATE UNIQUE INDEX idx_types__m1 ON types (t_int16, t_int32);

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
	typeId     {{.MySQLInt}},
	t_duration TIME,
	t_datetime TIMESTAMP
);

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

	FOREIGN KEY (ref_num, ref_type) REFERENCES account (acc_num, acc_type)
);
CREATE INDEX idx_sub_account__m1 ON sub_account (ref_num, ref_type);

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

CREATE TABLE `user` (
	user_id    {{.MySQLInt}} PRIMARY KEY,
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

CREATE TABLE user_address (
	user_id    {{.MySQLInt}} REFERENCES `user`(user_id),
	address_id {{.MySQLInt}} REFERENCES address(address_id),

	PRIMARY KEY (user_id, address_id)
);

INSERT INTO sex (id, name)
	VALUES(0, 'female');
INSERT INTO sex (id, name)
	VALUES(1, 'male');

INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool)
	VALUES(1, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 'A', 'Z', TRUE);

INSERT INTO times (typeId, t_duration, t_datetime)
	VALUES(1, '5:3:12', '2009-11-10T23:00:00Z');
INSERT INTO times (typeId, t_duration, t_datetime)
	VALUES(2, NULL, NULL);
COMMIT;
