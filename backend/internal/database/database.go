package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Migrate() error {
	migrationFiles := []string{
		"001_initial_schema.sql",
	}

	for _, file := range migrationFiles {
		path := filepath.Join("migrations", file)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
