package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      date CHAR(8) NOT NULL DEFAULT '',
      title VARCHAR(255) NOT NULL DEFAULT '',
      comment TEXT NOT NULL DEFAULT '',
      repeat VARCHAR(128) NOT NULL DEFAULT ''
   );`
	Name = "scheduler.db"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbFile string) (*Database, error) {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("can't open DB: %w", err)
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			return nil, fmt.Errorf("can't create DB: %w", err)
		}
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
