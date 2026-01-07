package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/enneket/rednote-extract/tool/crawler"
	"github.com/joho/godotenv"
)

var (
	keyword = flag.String("k", "", "keyword to extract")
)

func main() {
	flag.Parse()
	if *keyword == "" {
		flag.Usage()
		return
	}

	// 加载 ENV
	godotenv.Load()

	cookie := os.Getenv("cookie")

	// 创建爬虫实例
	crawler, err := crawler.NewCrawler(cookie)
	if err != nil {
		panic(err)
	}

	// 搜索帖子
	notes, err := crawler.Search(*keyword)
	if err != nil {
		panic(err)
	}
	for _, note := range notes {
		fmt.Printf("NoteID: %s, Title: %s, Author: %s, LikeCount: %d, CommentCount: %d\n",
			note.NoteID, note.Title, note.AuthorName, note.LikeCount, note.CommentCount)
	}
}
