package image

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	CreateTable() error
	SaveImage(data []byte) error
	GetImage(db *sql.DB) (*Image, error)
}

type Image struct {
	ID   int
	Data []byte
}

type SQLiteStore struct {
	DB *sql.DB
}

func NewStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("データベースのオープンに失敗: %w", err)
	}

	return &SQLiteStore{DB: db}, nil
}

func (s *SQLiteStore) CreateTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB
		)
	`)
	if err != nil {
		return fmt.Errorf("テーブルの作成に失敗: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SaveImage(data []byte) error {
	_, err := s.DB.Exec("INSERT INTO images (data) VALUES (?)", data)
	if err != nil {
		return fmt.Errorf("画像の保存に失敗: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetImage(db *sql.DB) (*Image, error) {
	var img Image
	err := s.DB.QueryRow("SELECT id, data FROM images ORDER BY id DESC LIMIT 1").Scan(&img.ID, &img.Data)
	if err != nil {
		return nil, fmt.Errorf("画像の取得に失敗: %w", err)
	}
	return &img, nil
}
