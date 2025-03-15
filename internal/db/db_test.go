package db

import (
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	// テスト用の一時データベースファイルを作成
	tmpDB := "test.db"
	defer os.Remove(tmpDB)

	// データベース接続
	store, err := NewStore(tmpDB)
	if err != nil {
		t.Fatalf("データベースの接続に失敗: %v", err)
	}
	defer store.Close()

	// テーブル作成
	err = store.CreateTable()
	if err != nil {
		t.Fatalf("テーブルの作成に失敗: %v", err)
	}

	// テスト用の画像データ
	testData := []byte("test image data")

	// 画像の保存
	err = store.SaveImage(testData)
	if err != nil {
		t.Fatalf("画像の保存に失敗: %v", err)
	}

	// 画像の取得
	img, err := store.GetImage(store.DB)
	if err != nil {
		t.Fatalf("画像の取得に失敗: %v", err)
	}

	// 取得したデータの検証
	if string(img.Data) != string(testData) {
		t.Errorf("取得したデータが一致しません。\n期待値: %s\n実際の値: %s", testData, img.Data)
	}
}
