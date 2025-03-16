package image

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
)

type Image struct {
	ID   int
	URL  string
	Name string
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
			url TEXT NOT NULL,
			name TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("テーブルの作成に失敗: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListImages() ([]Image, error) {
	rows, err := s.DB.Query("SELECT id, url, name FROM images")
	if err != nil {
		return nil, fmt.Errorf("画像の一覧取得に失敗: %w", err)
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.URL, &img.Name); err != nil {
			return nil, fmt.Errorf("画像データの読み取りに失敗: %w", err)
		}
		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("画像の一覧取得中にエラーが発生: %w", err)
	}

	return images, nil
}

func (s *SQLiteStore) SaveImage(name, url string) error {
	_, err := s.DB.Exec("INSERT INTO images (url, name) VALUES (?, ?)", url, name)
	if err != nil {
		return fmt.Errorf("画像の保存に失敗: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetImage() (*Image, error) {
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
	err = s.DB.QueryRow("SELECT id, url, name FROM images LIMIT 1 OFFSET ?", offset).Scan(&img.ID, &img.URL, &img.Name)
	if err != nil {
		return nil, fmt.Errorf("画像の取得に失敗: %w", err)
	}
	return &img, nil
}

func (s *SQLiteStore) GetImageByName(name string) (*Image, error) {
	var img Image
	err := s.DB.QueryRow("SELECT id, url, name FROM images WHERE name = ?", name).Scan(&img.ID, &img.URL, &img.Name)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("指定された名前の画像が見つかりません: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("画像の取得に失敗: %w", err)
	}
	return &img, nil
}
