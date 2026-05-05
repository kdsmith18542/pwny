package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
	path string
}

func Open(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d := &Database{DB: db, path: path}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	slog.Info("database opened", "path", path)
	return d, nil
}

func (d *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS workspaces (
		id          TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		description TEXT,
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS hosts (
		id          TEXT PRIMARY KEY,
		workspace_id TEXT REFERENCES workspaces(id),
		address     TEXT NOT NULL,
		mac         TEXT,
		os_name     TEXT,
		os_flavor   TEXT,
		os_sp       TEXT,
		arch        TEXT,
		purpose     TEXT,
		info        TEXT,
		state       TEXT DEFAULT 'alive',
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS services (
		id          TEXT PRIMARY KEY,
		host_id     TEXT REFERENCES hosts(id),
		port        INTEGER NOT NULL,
		proto       TEXT NOT NULL,
		name        TEXT,
		info        TEXT,
		state       TEXT DEFAULT 'open',
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS credentials (
		id           TEXT PRIMARY KEY,
		workspace_id TEXT REFERENCES workspaces(id),
		host_id      TEXT REFERENCES hosts(id),
		username     TEXT NOT NULL,
		password_enc TEXT,
		hash_enc     TEXT,
		type         TEXT,
		module       TEXT,
		created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS notes (
		id          TEXT PRIMARY KEY,
		workspace_id TEXT REFERENCES workspaces(id),
		host_id     TEXT REFERENCES hosts(id),
		title       TEXT NOT NULL,
		content     TEXT,
		ntype       TEXT DEFAULT 'note',
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS event_log (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id TEXT,
		level       TEXT NOT NULL,
		module      TEXT,
		message     TEXT NOT NULL,
		detail      TEXT,
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.Exec(schema)
	return err
}

func (d *Database) Close() error {
	slog.Info("database closing", "path", d.path)
	return d.DB.Close()
}
