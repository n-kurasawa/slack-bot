package image

import (
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	tmpDB := "test.db"
	defer os.Remove(tmpDB)

	store, err := NewStore(tmpDB)
	if err != nil {
		t.Fatal(err)
	}
	defer store.DB.Close()

	if err := store.CreateTable(); err != nil {
		t.Fatal(err)
	}

	// 複数の画像URLを保存
	testURLs := []string{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
	}

	for _, url := range testURLs {
		if err := store.SaveImage(url); err != nil {
			t.Fatal(err)
		}
	}

	// 画像が取得できることを確認
	img, err := store.GetImage()
	if err != nil {
		t.Fatal(err)
	}

	// 取得した画像URLが保存した画像URLのいずれかと一致することを確認
	found := false
	for _, url := range testURLs {
		if img.URL == url {
			found = true
			break
		}
	}

	if !found {
		t.Error("取得した画像URLが保存した画像URLと一致しません")
	}

	// 画像一覧を取得して確認
	images, err := store.ListImages()
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != len(testURLs) {
		t.Errorf("保存した画像の数が一致しません: want %d, got %d", len(testURLs), len(images))
	}
}
