BEGIN TRANSACTION;

CREATE TABLE types (
    field_int INTEGER,
    field_float FLOAT,
    field_text TEXT,
    field_blob BLOB,
    field_bool BOOLEAN);

CREATE TABLE _types (id TEXT PRIMARY KEY,
    field_int TEXT,
    field_float TEXT,
    field_text TEXT,
    field_blob TEXT,
    field_bool TEXT);

CREATE TABLE default_value (id INTEGER PRIMARY KEY,
    def_int INTEGER DEFAULT 55,
    def_float FLOAT DEFAULT 10.1,
    def_text TEXT DEFAULT 'string',
    def_bool BOOLEAN DEFAULT 0);

CREATE TABLE _default_value (id TEXT PRIMARY KEY,
    def_int TEXT,
    def_float TEXT,
    def_text TEXT,
    def_bool TEXT);

COMMIT;
