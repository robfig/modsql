BEGIN TRANSACTION;

CREATE TABLE types (
    t_int INTEGER,
    t_float FLOAT,
    t_text TEXT,
    t_blob BLOB,
    t_bool BOOLEAN);

CREATE TABLE _types (id TEXT PRIMARY KEY,
    t_int TEXT,
    t_float TEXT,
    t_text TEXT,
    t_blob TEXT,
    t_bool TEXT);

CREATE TABLE default_value (id INTEGER PRIMARY KEY,
    d_int INTEGER DEFAULT 55,
    d_float FLOAT DEFAULT 10.2,
    d_text TEXT DEFAULT 'string',
    d_bool BOOLEAN DEFAULT 0);

CREATE TABLE _default_value (id TEXT PRIMARY KEY,
    d_int TEXT,
    d_float TEXT,
    d_text TEXT,
    d_bool TEXT);

COMMIT;
BEGIN TRANSACTION;

INSERT INTO "types" (t_int, t_float, t_text, t_blob, t_bool) VALUES(1, 1.1, 'one', 'one', 1);
INSERT INTO "types" (t_int, t_float, t_text, t_blob, t_bool) VALUES(2, 2.2, 'two', 'two', 0);

INSERT INTO "default_value" (id, d_int, d_float, d_text, d_bool) VALUES(1, 10, 10.1, 'foo', 1);

COMMIT;
