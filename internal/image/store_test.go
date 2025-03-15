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

	// 複数の画像を保存
	testData := [][]byte{
		[]byte("test data 1"),
		[]byte("test data 2"),
		[]byte("test data 3"),
	}

	for _, data := range testData {
		if err := store.SaveImage(data); err != nil {
			t.Fatal(err)
		}
	}

	// 画像が取得できることを確認
	img, err := store.GetImage()
	if err != nil {
		t.Fatal(err)
	}

	// 取得した画像データが保存した画像データのいずれかと一致することを確認
	found := false
	for _, data := range testData {
		if string(img.Data) == string(data) {
			found = true
			break
		}
	}

	if !found {
		t.Error("取得した画像データが保存した画像データと一致しません")
	}
}
