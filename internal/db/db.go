package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/n-kurasawa/slack-bot/internal/handler"
)

type Store struct {
	DB *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("データベースの接続に失敗: %w", err)
	}

	return &Store{DB: db}, nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}

func (s *Store) CreateTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("テーブルの作成に失敗: %w", err)
	}
	return nil
}

func (s *Store) SaveImage(data []byte) error {
	_, err := s.DB.Exec("INSERT INTO images (data) VALUES (?)", data)
	if err != nil {
		return fmt.Errorf("画像の保存に失敗: %w", err)
	}
	return nil
}

func (s *Store) GetImage(db *sql.DB) (*handler.Image, error) {
	var img handler.Image
	err := s.DB.QueryRow("SELECT id, data FROM images ORDER BY id DESC LIMIT 1").Scan(&img.ID, &img.Data)
	if err != nil {
		return nil, fmt.Errorf("画像の取得に失敗: %w", err)
	}
	return &img, nil
}
