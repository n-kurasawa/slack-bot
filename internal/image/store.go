package image

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
)

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
	// 画像の総数を取得
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM images").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("画像数の取得に失敗: %w", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("画像が登録されていません")
	}

	// ランダムなオフセットを生成
	offset := rand.Intn(count)

	var img Image
	err = s.DB.QueryRow("SELECT id, data FROM images LIMIT 1 OFFSET ?", offset).Scan(&img.ID, &img.Data)
	if err != nil {
		return nil, fmt.Errorf("画像の取得に失敗: %w", err)
	}
	return &img, nil
}
