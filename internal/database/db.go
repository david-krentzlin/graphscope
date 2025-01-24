package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE types (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		kind TEXT NOT NULL
	);

	CREATE TABLE directives (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		type_id INTEGER,
		FOREIGN KEY(type_id) REFERENCES types(id)
	);

	CREATE TABLE directive_arguments (
		id INTEGER PRIMARY KEY,
		directive_id INTEGER,
		name TEXT NOT NULL,
		value_type TEXT NOT NULL,
		default_value TEXT,
		FOREIGN KEY(directive_id) REFERENCES directives(id)
	);

	CREATE TABLE directive_locations (
		directive_id INTEGER,
		location TEXT NOT NULL,
		FOREIGN KEY(directive_id) REFERENCES directives(id)
	);`

	_, err := db.Exec(schema)
	return err
}
