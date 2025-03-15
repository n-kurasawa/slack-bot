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
	list := flag.Bool("list", false, "画像の一覧を表示")
	flag.Parse()

	store, err := image.NewStore(*dbPath)
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
			fmt.Printf("ID: %d, サイズ: %d bytes\n", img.ID, len(img.Data))
		}
		return
	}

	if *imagePath == "" {
		fmt.Println("使用方法:")
		fmt.Println("  画像の登録: -image <画像ファイルのパス>")
		fmt.Println("  画像の一覧: -list")
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
