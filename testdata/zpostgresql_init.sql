// +build PostgreSQL
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE sex (
	id   smallint PRIMARY KEY,
	name text
);

CREATE TABLE types (
	t_int     {{.PostgreInt}} PRIMARY KEY,
	t_int8    smallint,
	t_int16   smallint,
	t_int32   integer,
	t_int64   bigint,
	t_float32 real,
	t_float64 double precision,
	t_string  text UNIQUE,
	t_binary  bytea,
	t_byte    character,
	t_rune    character varying(4),
	t_bool    boolean,

	UNIQUE (t_float32, t_float64)
);
CREATE UNIQUE INDEX idx_types_t_float64 ON types (t_float64);
CREATE INDEX idx_types_t_rune ON types (t_rune);
CREATE UNIQUE INDEX idx_types__m1 ON types (t_int16, t_int32);

CREATE TABLE default_value (
	id        {{.PostgreInt}} PRIMARY KEY,
	d_int8    smallint DEFAULT 55,
	d_float32 real DEFAULT 10.2,
	d_string  text,
	d_binary  bytea,
	d_byte    character DEFAULT 'b',
	d_rune    character varying(4) DEFAULT 'r',
	d_bool    boolean DEFAULT FALSE
);

CREATE TABLE times (
	typeId     {{.PostgreInt}},
	t_duration time without time zone,
	t_datetime timestamp without time zone
);

CREATE TABLE account (
	acc_num   {{.PostgreInt}},
	acc_type  {{.PostgreInt}},
	acc_descr text,

	PRIMARY KEY (acc_num, acc_type)
);

CREATE TABLE sub_account (
	sub_acc   {{.PostgreInt}} PRIMARY KEY,
	ref_num   {{.PostgreInt}},
	ref_type  {{.PostgreInt}},
	sub_descr text,

	FOREIGN KEY (ref_num, ref_type) REFERENCES account (acc_num, acc_type)
);
CREATE INDEX idx_sub_account__m1 ON sub_account (ref_num, ref_type);

CREATE TABLE catalog (
	catalog_id  {{.PostgreInt}} PRIMARY KEY,
	name        text,
	description text,
	price       real
);

CREATE TABLE magazine (
	catalog_id {{.PostgreInt}} PRIMARY KEY REFERENCES catalog(catalog_id),
	page_count text
);

CREATE TABLE mp3 (
	catalog_id {{.PostgreInt}} PRIMARY KEY REFERENCES catalog(catalog_id),
	size       {{.PostgreInt}},
	length     real,
	filename   text
);

CREATE TABLE book (
	book_id {{.PostgreInt}} PRIMARY KEY,
	title   text,
	author  text
);

CREATE TABLE chapter (
	chapter_id {{.PostgreInt}} PRIMARY KEY,
	title      text,
	book_fk    {{.PostgreInt}} REFERENCES book(book_id)
);

CREATE TABLE "user" (
	user_id    {{.PostgreInt}} PRIMARY KEY,
	first_name text,
	last_name  text
);

CREATE TABLE address (
	address_id {{.PostgreInt}} PRIMARY KEY,
	street     text,
	city       text,
	state      text,
	post_code  text
);

CREATE TABLE user_address (
	user_id    {{.PostgreInt}} REFERENCES "user"(user_id),
	address_id {{.PostgreInt}} REFERENCES address(address_id),

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
