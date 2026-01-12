package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enneket/rednote-extract/agent"
	"github.com/enneket/rednote-extract/config"
	"github.com/enneket/rednote-extract/models"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// Get today's date folder
	today := time.Now().Format("2006-01-02")
	dataDir := filepath.Join("data", today)

	// Read all JSON files from the data folder for today
	noteInputs, err := readNotesFromFolder(dataDir)
	if err != nil {
		log.Fatalf("Failed to read notes from folder: %v", err)
	}

	if len(noteInputs) == 0 {
		log.Println("No notes found in data folder. Using default note.")
		panic("No notes found in data folder.")
	}

	fmt.Printf("=== 加载了 %d 篇原始笔记 ===\n", len(noteInputs))
	for i, input := range noteInputs {
		fmt.Printf("第 %d 篇笔记:\n", i+1)
		fmt.Printf("  标题: %s\n", input.Title)
		fmt.Printf("  内容: %.100s...\n", input.Content) // Show first 100 chars
		fmt.Printf("  评论数: %d\n\n", len(input.Comments))
	}

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

			// Read the JSON file
			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", filePath, err)
				continue
			}

			// Parse JSON data based on file type
			if strings.Contains(strings.ToLower(file.Name()), "contents") {
				// Process content files
				var contents []ContentData
				if err := json.Unmarshal(data, &contents); err != nil {
					log.Printf("Failed to parse content file %s: %v", filePath, err)
					continue
				}

				// Store content by note_id
				for _, content := range contents {
					contentMap[content.NoteID] = content
				}
			} else if strings.Contains(strings.ToLower(file.Name()), "comments") {
				// Process comment files
				var comments []CommentData
				if err := json.Unmarshal(data, &comments); err != nil {
					log.Printf("Failed to parse comment file %s: %v", filePath, err)
					continue
				}

				// Group comments by note_id
				for _, comment := range comments {
					commentMap[comment.NoteID] = append(commentMap[comment.NoteID], comment.Content)
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
