// +build sqlite
// MACHINE GENERATED BY ModSQL (github.com/kless/modsql); DO NOT EDIT

BEGIN;
CREATE TABLE types (
	t_int     INTEGER PRIMARY KEY,
	t_int8    INTEGER,
	t_int16   INTEGER,
	t_int32   INTEGER,
	t_int64   INTEGER,
	t_float32 REAL,
	t_float64 REAL,
	t_string  VARCHAR(32) UNIQUE,
	t_binary  BLOB,
	t_byte    TEXT,
	t_rune    TEXT,
	t_bool    BOOL,

	UNIQUE (t_float32, t_float64)
);
CREATE UNIQUE INDEX ix_types_t_float64 ON types (t_float64);
CREATE INDEX ix_types_t_rune ON types (t_rune);
CREATE INDEX ix_types__m1 ON types (t_int8, t_float32);
CREATE UNIQUE INDEX ix_types__m2 ON types (t_int16, t_int32);

CREATE TABLE default_value (
	id        INTEGER PRIMARY KEY,
	d_int8    INTEGER DEFAULT 55,
	d_float32 REAL DEFAULT 10.2,
	d_string  TEXT,
	d_binary  BLOB,
	d_byte    TEXT DEFAULT 'b',
	d_rune    TEXT DEFAULT 'r',
	d_bool    BOOL DEFAULT 0,
	d_findex  INTEGER
);

CREATE TABLE times (
	typeId     INTEGER REFERENCES types(t_int),
	t_duration INTEGER,
	t_datetime TEXT,

	PRIMARY KEY (t_duration, t_datetime)
);
CREATE INDEX ix_times_t_datetime ON times (t_datetime);

INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool) VALUES(1, 8, 16, 32, 64, 1.32, 1.64, 'one', '12', 'A', 'Z', 1);
INSERT INTO default_value (id, d_int8, d_float32, d_string, d_binary, d_byte, d_rune, d_bool, d_findex) VALUES(1, 10, 10.1, 'foo', '12', 'a', 'z', 0, 1);
INSERT INTO times (typeId, t_duration, t_datetime) VALUES(1, '5:3:12', '2012-09-22 19:29:08');
COMMIT;
