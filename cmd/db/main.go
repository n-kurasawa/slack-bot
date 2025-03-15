package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/n-kurasawa/slack-bot/internal/image"
)

func main() {
	dbPath := flag.String("db", "images.db", "データベースファイルのパス")
	imagePath := flag.String("image", "", "画像ファイルのパス")
	flag.Parse()

	store, err := image.NewStore(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	if *imagePath == "" {
		fmt.Println("画像ファイルのパスを指定してください")
		os.Exit(1)
	}

	data, err := os.ReadFile(*imagePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.SaveImage(data); err != nil {
		log.Fatal(err)
	}

	fmt.Println("画像を保存しました")
}
