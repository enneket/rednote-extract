package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/enneket/rednote-extract/tool/crawler/internal/model"
)

// JSONStore JSON存储实现
type JSONStore struct {
	dirPath  string
	notes    map[string]*model.Note
	comments map[string][]*model.Comment
	mutex    sync.RWMutex
}

// NewJSONStore 创建JSON存储实例
func NewJSONStore(config map[string]interface{}) (Store, error) {
	dirPath := "./data"
	if path, ok := config["dir"].(string); ok {
		dirPath = path
	}

	// 创建数据目录
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &JSONStore{
		dirPath:  dirPath,
		notes:    make(map[string]*model.Note),
		comments: make(map[string][]*model.Comment),
	}, nil
}

// SaveNote 保存帖子
func (s *JSONStore) SaveNote(note *model.Note) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.notes[note.NoteID] = note

	// 保存到文件
	filePath := filepath.Join(s.dirPath, fmt.Sprintf("note_%s.json", note.NoteID))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create note file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(note); err != nil {
		return fmt.Errorf("failed to encode note: %w", err)
	}

	return nil
}

// SaveComment 保存评论
func (s *JSONStore) SaveComment(comment *model.Comment) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.comments[comment.NoteID] = append(s.comments[comment.NoteID], comment)

	// 保存到文件
	filePath := filepath.Join(s.dirPath, fmt.Sprintf("comments_%s.json", comment.NoteID))

	// 读取现有评论
	var existingComments []*model.Comment
	existingFile, err := os.Open(filePath)
	if err == nil {
		defer existingFile.Close()
		if err := json.NewDecoder(existingFile).Decode(&existingComments); err == nil {
			// 去重
			existingMap := make(map[string]bool)
			for _, c := range existingComments {
				existingMap[c.CommentID] = true
			}

			if !existingMap[comment.CommentID] {
				existingComments = append(existingComments, comment)
			}
		} else {
			existingComments = []*model.Comment{comment}
		}
	} else {
		existingComments = []*model.Comment{comment}
	}

	// 写入文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create comments file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingComments); err != nil {
		return fmt.Errorf("failed to encode comments: %w", err)
	}

	return nil
}

// SaveMedia 保存媒体文件
func (s *JSONStore) SaveMedia(noteID, mediaType, fileName string, content []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 创建媒体目录
	mediaDir := filepath.Join(s.dirPath, "media", noteID, mediaType)
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// 保存文件
	filePath := filepath.Join(mediaDir, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write media file: %w", err)
	}

	return nil
}

// GetNoteByID 根据ID获取帖子
func (s *JSONStore) GetNoteByID(noteID string) (*model.Note, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	note, ok := s.notes[noteID]
	if !ok {
		// 从文件读取
		filePath := filepath.Join(s.dirPath, fmt.Sprintf("note_%s.json", noteID))
		file, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to open note file: %w", err)
		}
		defer file.Close()

		var noteFromFile model.Note
		if err := json.NewDecoder(file).Decode(&noteFromFile); err != nil {
			return nil, fmt.Errorf("failed to decode note: %w", err)
		}

		s.notes[noteID] = &noteFromFile
		return &noteFromFile, nil
	}

	return note, nil
}

// GetCommentsByNoteID 根据帖子ID获取评论
func (s *JSONStore) GetCommentsByNoteID(noteID string) ([]*model.Comment, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	comments, ok := s.comments[noteID]
	if !ok || len(comments) == 0 {
		// 从文件读取
		filePath := filepath.Join(s.dirPath, fmt.Sprintf("comments_%s.json", noteID))
		file, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return []*model.Comment{}, nil
			}
			return nil, fmt.Errorf("failed to open comments file: %w", err)
		}
		defer file.Close()

		var commentsFromFile []*model.Comment
		if err := json.NewDecoder(file).Decode(&commentsFromFile); err != nil {
			return nil, fmt.Errorf("failed to decode comments: %w", err)
		}

		s.comments[noteID] = commentsFromFile
		return commentsFromFile, nil
	}

	return comments, nil
}

// Close 关闭存储
func (s *JSONStore) Close() error {
	// JSON存储不需要特殊关闭操作
	return nil
}
