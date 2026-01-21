package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enneket/rednote-extract/internal/models"
)

// Struct to match the content JSON format
type ContentData struct {
	NoteID string `json:"note_id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
}

type SavedCommentItem struct {
	Content string `json:"content"`
}

type SavedComments struct {
	NoteID   string             `json:"note_id"`
	Comments []SavedCommentItem `json:"comments"`
}

// ReadNotesFromFolder reads all JSON files from the specified folder and converts them to []*models.NoteInput
func ReadNotesFromFolder(folderPath string, maxNotes int) ([]*models.NoteInput, error) {
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
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			continue
		}
		
		filePath := filepath.Join(folderPath, file.Name())
		log.Printf("Processing file: %s", filePath)

		if err := processFile(filePath, file.Name(), contentMap, commentMap); err != nil {
			log.Printf("Failed to process file %s: %v", filePath, err)
		}
	}

	// Combine content and comments into NoteInput objects
	var noteInputs []*models.NoteInput
	count := 0
	// Process content in order of note_id to maintain consistency (map iteration is random, but acceptable here)
	for noteID, content := range contentMap {
		if maxNotes > 0 && count >= maxNotes {
			log.Printf("Reached maximum note limit of %d, stopping processing", maxNotes)
			break
		}

		noteInput := &models.NoteInput{
			Title:    content.Title,
			Content:  content.Desc,
			Comments: commentMap[noteID],
		}

		noteInputs = append(noteInputs, noteInput)
		count++
	}
	log.Printf("Loaded %d notes (max %d)", len(noteInputs), maxNotes)
	return noteInputs, nil
}

func processFile(filePath, fileName string, contentMap map[string]ContentData, commentMap map[string][]string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	lowerName := strings.ToLower(fileName)

	if strings.Contains(lowerName, "contents") || strings.Contains(lowerName, "notes") {
		for {
			var content ContentData
			if err := decoder.Decode(&content); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding content: %w", err)
			}
			if content.NoteID != "" {
				contentMap[content.NoteID] = content
			}
		}
	} else if strings.Contains(lowerName, "comments") {
		for {
			var saved SavedComments
			if err := decoder.Decode(&saved); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding comments: %w", err)
			}
			for _, c := range saved.Comments {
				commentMap[saved.NoteID] = append(commentMap[saved.NoteID], c.Content)
			}
		}
	}
	return nil
}

// SaveGeneratedNote saves the generated note to a file
func SaveGeneratedNote(note *models.GeneratedNote) error {
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
