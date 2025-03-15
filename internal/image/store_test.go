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

	testData := []byte("test data")
	if err := store.SaveImage(testData); err != nil {
		t.Fatal(err)
	}

	img, err := store.GetImage(store.DB)
	if err != nil {
		t.Fatal(err)
	}

	if string(img.Data) != string(testData) {
		t.Errorf("want %s, got %s", testData, img.Data)
	}
}
