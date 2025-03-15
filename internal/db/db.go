package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Image struct {
	ID   int64
	Data []byte
}

func OpenDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL
		)
	`)
	return err
}

func SaveImage(db *sql.DB, data []byte) error {
	_, err := db.Exec("INSERT INTO images (data) VALUES (?)", data)
	return err
}

func GetImage(db *sql.DB) (*Image, error) {
	var img Image
	err := db.QueryRow("SELECT id, data FROM images ORDER BY id DESC LIMIT 1").Scan(&img.ID, &img.Data)
	if err != nil {
		return nil, err
	}
	return &img, nil
}
