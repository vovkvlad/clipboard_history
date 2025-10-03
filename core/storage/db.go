package storage

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, error) {
	// TODO: Get db path based on the OS
	db, err := sql.Open("sqlite3", path.Join(".", "tmp", "clipboard_history.db"))

	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}
