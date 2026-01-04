package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/enneket/rednote-extract/browser"
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
	fmt.Println(*keyword)
	// 根据关键词在小红书进行搜索
	notes, err := browser.SearchRednote(context.Background(), *keyword)
	if err != nil {
		time.Sleep(2 * time.Minute)
		fmt.Println("SearchRednote failed:", err)
		return
	}
	fmt.Println("SearchRednote success:", notes)

	// 根据搜索结果查看每个笔记的内容和评论
	time.Sleep(2 * time.Minute)
}
