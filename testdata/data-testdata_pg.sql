// +build postgresql
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE types (
	t_int     {{.PostgreInt}} PRIMARY KEY,
	t_int8    smallint,
	t_int16   smallint,
	t_int32   integer,
	t_int64   bigint,
	t_float32 real,
	t_float64 double precision,
	t_string  text,
	t_binary  bytea,
	t_byte    character,
	t_rune    character varying(4),
	t_bool    boolean
);

CREATE TABLE _types (
	lang      VARCHAR(32) PRIMARY KEY,
	t_int     TEXT,
	t_int8    TEXT,
	t_int16   TEXT,
	t_int32   TEXT,
	t_int64   TEXT,
	t_float32 TEXT,
	t_float64 TEXT,
	t_string  TEXT,
	t_binary  TEXT,
	t_byte    TEXT,
	t_rune    TEXT,
	t_bool    TEXT
);

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

CREATE TABLE _default_value (
	lang      VARCHAR(32) PRIMARY KEY,
	id        TEXT,
	d_int8    TEXT,
	d_float32 TEXT,
	d_string  TEXT,
	d_binary  TEXT,
	d_byte    TEXT,
	d_rune    TEXT,
	d_bool    TEXT
);

CREATE TABLE times (
	t_duration time without time zone,
	t_datetime timestamp without time zone
);

CREATE TABLE _times (
	lang       VARCHAR(32) PRIMARY KEY,
	t_duration TEXT,
	t_datetime TEXT
);

INSERT INTO _types (lang, t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool) VALUES('en', 'int', 'integer 8', 'integer 16', 'integer 32', 'integer 64', 'float 32', 'float 64', 'string', 'binary', 'byte', 'rune', 'boolean');
INSERT INTO _default_value (lang, id, d_int8, d_float32, d_string, d_binary, d_byte, d_rune, d_bool) VALUES('en', 'id', 'integer 8', 'float 32', 'string', 'binary', 'byte', 'rune', 'boolean');
INSERT INTO _times (lang, t_duration, t_datetime) VALUES('en', 'duration', 'datetime');
INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool) VALUES(1, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 'A', 'Z', TRUE);
INSERT INTO default_value (id, d_int8, d_float32, d_string, d_binary, d_byte, d_rune, d_bool) VALUES(1, 10, 10.1, 'foo', '12', 'a', 'z', FALSE);
INSERT INTO times (t_duration, t_datetime) VALUES('5:3:12', '2012-08-09 09:52:39');
COMMIT;
