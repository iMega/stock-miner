package storage

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Option func(b *Storage)

type Storage struct {
	db *sql.DB
}

func WithSqllite(db *sql.DB) Option {
	return func(s *Storage) {
		s.db = db
	}
}

func New(opts ...Option) *Storage {
	s := &Storage{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func CreateDatabase(name string) error {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(name, os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return err
	}

	query := `CREATE TABLE price (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol VARCHAR(64) NULL,
        price DECIMAL(10,5) NULL,
        create_at DATETIME NULL
    )`

	if _, err = db.Exec(query); err != nil {
		return err
	}

	return db.Close()
}
