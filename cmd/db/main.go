package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/n-kurasawa/slack-bot/internal/db"
)

func main() {
	var (
		dbPath    = flag.String("db", "images.db", "データベースファイルのパス")
		imagePath = flag.String("image", "", "登録する画像ファイルのパス")
	)
	flag.Parse()

	database, err := db.OpenDB(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.CreateTable(database); err != nil {
		log.Fatal(err)
	}

	if *imagePath != "" {
		data, err := os.ReadFile(*imagePath)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.SaveImage(database, data); err != nil {
			log.Fatal(err)
		}
		fmt.Println("画像を登録しました")
	}
}
