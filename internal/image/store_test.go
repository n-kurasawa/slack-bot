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

	// テストデータの準備
	testData := []struct {
		url  string
		name string
	}{
		{url: "https://example.com/image1.jpg", name: "cat"},
		{url: "https://example.com/image2.jpg", name: "dog"},
		{url: "https://example.com/image3.jpg", name: "bird"},
	}

	// 画像の保存
	for _, data := range testData {
		if err := store.SaveImage(data.url + " " + data.name); err != nil {
			t.Fatal(err)
		}
	}

	// 名前なしで画像を保存しようとした場合のエラーをテスト
	t.Run("名前なしでの画像保存", func(t *testing.T) {
		err := store.SaveImage("https://example.com/image4.jpg")
		if err == nil {
			t.Error("名前なしでの画像保存がエラーを返しませんでした")
		}
	})

	// ランダムな画像の取得をテスト
	t.Run("ランダムな画像の取得", func(t *testing.T) {
		img, err := store.GetImage()
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, data := range testData {
			if img.URL == data.url && img.Name == data.name {
				found = true
				break
			}
		}

		if !found {
			t.Error("取得した画像が保存した画像と一致しません")
		}
	})

	// 名前による画像の取得をテスト
	t.Run("名前による画像の取得", func(t *testing.T) {
		img, err := store.GetImageByName("cat")
		if err != nil {
			t.Fatal(err)
		}

		if img.URL != testData[0].url {
			t.Errorf("画像URLが一致しません: want %s, got %s", testData[0].url, img.URL)
		}

		if img.Name != testData[0].name {
			t.Errorf("画像の名前が一致しません: want %s, got %s", testData[0].name, img.Name)
		}
	})

	// 存在しない名前による画像の取得をテスト
	t.Run("存在しない名前による画像の取得", func(t *testing.T) {
		_, err := store.GetImageByName("nonexistent")
		if err == nil {
			t.Error("存在しない名前による画像の取得がエラーを返しませんでした")
		}
	})

	// 画像一覧を取得して確認
	t.Run("画像一覧の取得", func(t *testing.T) {
		images, err := store.ListImages()
		if err != nil {
			t.Fatal(err)
		}

		if len(images) != len(testData) {
			t.Errorf("保存した画像の数が一致しません: want %d, got %d", len(testData), len(images))
		}

		for _, img := range images {
			found := false
			for _, data := range testData {
				if img.URL == data.url && img.Name == data.name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("一覧に含まれる画像が期待したものと一致しません: %+v", img)
			}
		}
	})
}
