CREATE TABLE instance_params (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    sha256 TEXT NOT NULL COLLATE NOCASE,
    param blob NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS instance_params__sha256_index ON instance_params (sha256);

CREATE TABLE instances (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    created_at TIMESTAMP NOT NULL,
    param_id INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS instances__created_at_index ON instances (created_at);

CREATE INDEX IF NOT EXISTS instances__param_id_index ON instances (param_id);

---
CREATE TABLE sessions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    instance_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions__instance_id_index ON sessions (instance_id);

CREATE UNIQUE INDEX IF NOT EXISTS sessions__name_instance_id_index ON sessions (name, instance_id);

CREATE INDEX IF NOT EXISTS sessions__created_at_index ON sessions (created_at);

---
CREATE TABLE task_params (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    sha256 TEXT NOT NULL COLLATE NOCASE,
    param BLOB NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS task_params__sha256_index ON task_params (sha256);

---
CREATE TABLE tasks (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    param_id INTEGER NOT NULL,
    session_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    result BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS tasks__session_id_index ON tasks (session_id);

CREATE UNIQUE INDEX IF NOT EXISTS tasks__name_session_id_index ON tasks (name, session_id);

CREATE INDEX IF NOT EXISTS tasks__created_at_index ON tasks (created_at);

CREATE INDEX IF NOT EXISTS tasks__param_id_index ON tasks (param_id);

---
CREATE TABLE log_entries (
    kind INTEGER NOT NULL,
    assoc_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    lvl INTEGER NOT NULL,
    message TEXT NOT NULL COLLATE BINARY
);

CREATE INDEX IF NOT EXISTS log_entries__kind_assoc_id_index ON log_entries (kind, assoc_id);

CREATE INDEX IF NOT EXISTS log_entries__created_at_index ON log_entries (created_at);

CREATE INDEX IF NOT EXISTS log_entries__lvl_index ON log_entries (lvl);

---
CREATE TABLE file_metas(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL COLLATE BINARY,
    meta blob
);

CREATE UNIQUE INDEX IF NOT EXISTS file_metas__path ON file_metas(path);

---
CREATE TABLE file_task(
    meta_id integer NOT NULL,
    task_id INTEGER NOT NULL,
    UNIQUE (meta_id, task_id) ON CONFLICT FAIL
);

CREATE INDEX IF NOT EXISTS file_task__task_id_index ON file_task (task_id);

CREATE INDEX IF NOT EXISTS file_task__meta_id_index ON file_task (meta_id);