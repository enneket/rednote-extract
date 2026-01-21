package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enneket/rednote-extract/internal/agent"
	"github.com/enneket/rednote-extract/internal/config"
	"github.com/enneket/rednote-extract/internal/crawler/xhs"
	"github.com/enneket/rednote-extract/internal/models"
)

// Struct to match the content JSON format
type ContentData struct {
	NoteID string `json:"note_id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
}

// Struct to match the comments JSON format
type CommentData struct {
	NoteID  string `json:"note_id"`
	Content string `json:"content"`
}

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
	noteInputs, err := readNotesFromFolder(dataDir)
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
		if err := saveGeneratedNote(generatedNote); err != nil {
			log.Printf("Failed to save generated note: %v", err)
		} else {
			log.Println("Generated note saved successfully.")
		}
	}
}

// saveGeneratedNote saves the generated note to a file
func saveGeneratedNote(note *models.GeneratedNote) error {
	dir := filepath.Join("data", "output")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("generated_%s.json", time.Now().Format("20060102_150405"))
	path := filepath.Join(dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(note)
}

// readNotesFromFolder reads all JSON files from the specified folder and converts them to []*models.NoteInput
func readNotesFromFolder(folderPath string) ([]*models.NoteInput, error) {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Folder does not exist: %s", folderPath)
		return nil, nil
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %s: %w", folderPath, err)
	}

	// Map to store content by note_id
	contentMap := make(map[string]ContentData)

	// Map to store comments grouped by note_id
	commentMap := make(map[string][]string)

	// Process all JSON files
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			filePath := filepath.Join(folderPath, file.Name())
			log.Printf("Processing file: %s", filePath)

			f, err := os.Open(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", filePath, err)
				continue
			}
			defer f.Close()

			decoder := json.NewDecoder(f)

			if strings.Contains(strings.ToLower(file.Name()), "contents") || strings.Contains(strings.ToLower(file.Name()), "notes") {
				// Process content files
				for {
					var content ContentData
					if err := decoder.Decode(&content); err == io.EOF {
						break
					} else if err != nil {
						log.Printf("Error decoding content in %s: %v", filePath, err)
						// If decode fails, we might try to recover or just stop for this file
						// For now, let's stop to avoid infinite loops on bad data
						break
					}

					if content.NoteID != "" {
						contentMap[content.NoteID] = content
					}
				}
			} else if strings.Contains(strings.ToLower(file.Name()), "comments") {
				// Process comment files
				for {
					// Comments file structure: each object might be { "note_id": "...", "comments": [ ... ] } or individual comment?
					// In xhs/crawler.go:
					// data := map[string]interface{}{ "note_id": noteId, "comments": commentsRes.Comments }
					// So it's an object with a list of comments.

					type CommentsWrapper struct {
						NoteID   string      `json:"note_id"`
						Comments interface{} `json:"comments"` // Could be list of objects or strings?
						// models.Comment is likely a struct, let's check crawler usage.
						// client.GetNoteComments returns *models.NoteCommentsResponse
						// which has Comments []Comment
					}

					// Let's use a generic map to be safe first or look at models
					// But we defined CommentData struct as single comment:
					// type CommentData struct { NoteID string, Content string }
					// The previous code assumed CommentData structure.
					// BUT crawler saves:
					// map[string]interface{}{ "note_id": noteId, "comments": commentsRes.Comments }
					// where commentsRes.Comments is []models.Comment

					// So we need a struct that matches what Crawler saves.
					type SavedCommentItem struct {
						Content string `json:"content"`
					}
					type SavedComments struct {
						NoteID   string             `json:"note_id"`
						Comments []SavedCommentItem `json:"comments"`
					}

					var saved SavedComments
					if err := decoder.Decode(&saved); err == io.EOF {
						break
					} else if err != nil {
						log.Printf("Error decoding comments in %s: %v", filePath, err)
						break
					}

					// Extract comments content
					for _, c := range saved.Comments {
						commentMap[saved.NoteID] = append(commentMap[saved.NoteID], c.Content)
					}
				}
			}
		}
	}

	// Combine content and comments into NoteInput objects
	var noteInputs []*models.NoteInput
	maxNotes := 10
	count := 0
	// Process content in order of note_id to maintain consistency
	for noteID, content := range contentMap {
		if count >= maxNotes {
			log.Printf("Reached maximum note limit of %d, stopping processing", maxNotes)
			break
		}

		// Create a NoteInput with title, content, and associated comments
		noteInput := &models.NoteInput{
			Title:    content.Title,
			Content:  content.Desc,
			Comments: commentMap[noteID], // Associated comments for this note_id
		}

		noteInputs = append(noteInputs, noteInput)
		count++
	}
	log.Printf("Loaded %d notes (max %d)", len(noteInputs), maxNotes)
	return noteInputs, nil
}
