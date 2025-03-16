package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/n-kurasawa/slack-bot/internal/db"
)

func main() {
	dbPath := flag.String("db", "images.db", "データベースファイルのパス")
	imageURL := flag.String("url", "", "画像のURL")
	imageName := flag.String("name", "", "画像の名前")
	list := flag.Bool("list", false, "画像の一覧を表示")
	flag.Parse()

	store, err := db.NewStore(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	if *list {
		images, err := store.ListImages()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("登録されている画像の数: %d\n", len(images))
		for _, img := range images {
			if img.Name != "" {
				fmt.Printf("ID: %d, Name: %s, URL: %s\n", img.ID, img.Name, img.URL)
			} else {
				fmt.Printf("ID: %d, URL: %s\n", img.ID, img.URL)
			}
		}
		return
	}

	if *imageURL == "" || *imageName == "" {
		fmt.Println("使用方法:")
		fmt.Println("  画像の登録: -url <画像のURL> -name <画像の名前>")
		fmt.Println("  画像の一覧: -list")
		os.Exit(1)
	}

	if err := store.SaveImage(*imageName, *imageURL); err != nil {
		log.Fatal(err)
	}

	fmt.Println("画像を保存しました")
}
