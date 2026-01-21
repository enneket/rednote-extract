package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/enneket/rednote-extract/internal/agent"
	"github.com/enneket/rednote-extract/internal/config"
	"github.com/enneket/rednote-extract/internal/crawler/xhs"
	"github.com/enneket/rednote-extract/internal/storage"
)

func main() {
	keywords := flag.String("keywords", "", "Comma separated keywords to search")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *keywords != "" {
		cfg.Keywords = *keywords
	}

	ctx := context.Background()

	// 1. Start Crawler
	log.Println("=== 开始抓取 ===")
	crawler := xhs.NewCrawler()
	if err := crawler.Start(ctx); err != nil {
		log.Printf("Crawler finished with error (or just stopped): %v", err)
		// We continue even if crawler fails, to process existing data
	} else {
		log.Println("Crawler finished successfully.")
	}

	// 2. Read Data
	// Data is saved in data/xhs by the crawler
	dataDir := filepath.Join("data", "xhs")
	log.Printf("从 %s 读取数据...", dataDir)

	// Read all JSON files from the data folder
	noteInputs, err := storage.ReadNotesFromFolder(dataDir, cfg.MaxNotes)
	if err != nil {
		log.Fatalf("Failed to read notes from folder %s: %v", dataDir, err)
	}

	if len(noteInputs) == 0 {
		log.Println("No notes found in data folder.")
		return
	}

	fmt.Printf("=== 加载了 %d 篇原始笔记 ===\n", len(noteInputs))
	for i, input := range noteInputs {
		fmt.Printf("第 %d 篇笔记:\n", i+1)
		fmt.Printf("  标题: %s\n", input.Title)
		if len(input.Content) > 100 {
			fmt.Printf("  内容: %.100s...\n", input.Content)
		} else {
			fmt.Printf("  内容: %s\n", input.Content)
		}
		fmt.Printf("  评论数: %d\n\n", len(input.Comments))
	}

	// 3. Generate New Note
	if len(noteInputs) > 0 {
		noteAgent := agent.NewReactAgent(cfg)
		generatedNote, err := noteAgent.GenerateNote(ctx, noteInputs)
		if err != nil {
			log.Fatalf("Failed to generate note: %v", err)
		}

		fmt.Println("=== 生成的笔记 ===")
		fmt.Printf("标题: %s\n\n", generatedNote.Title)
		fmt.Println("正文:")
		fmt.Println(generatedNote.Content)
		fmt.Printf("\n话题标签: %s\n", generatedNote.Tags)

		fmt.Printf("\n字数统计: %d 字\n", len(generatedNote.Content))

		// 4. Save Generated Note
		if err := storage.SaveGeneratedNote(generatedNote); err != nil {
			log.Printf("Failed to save generated note: %v", err)
		} else {
			log.Println("Generated note saved successfully.")
		}
	}
}
