BEGIN TRANSACTION;

CREATE TABLE server (
    uuid TEXT,
    access_log TEXT,
    error_log TEXT,
    chroot TEXT DEFAULT "/var/www",
    pid_File TEXT,
    default_host INTEGER,
    name TEXT,
    port INTEGER);

CREATE TABLE _server (id TEXT PRIMARY KEY,
    uuid TEXT,
    access_log TEXT,
    error_log TEXT,
    chroot TEXT,
    pid_File TEXT,
    default_host TEXT,
    name TEXT,
    port TEXT);

CREATE TABLE host (id INTEGER PRIMARY KEY,
    server_id INTEGER,
    maintenance BOOLEAN DEFAULT 0,
    name TEXT,
    matching TEXT);

CREATE TABLE _host (id TEXT PRIMARY KEY,
    server_id TEXT,
    maintenance TEXT,
    name TEXT,
    matching TEXT);

COMMIT;
