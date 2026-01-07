package store

import (
	"github.com/enneket/rednote-extract/tool/crawler/internal/config"
	"github.com/enneket/rednote-extract/tool/crawler/internal/model"
)

// Store 数据存储接口
type Store interface {
	// SaveNote 保存帖子
	SaveNote(note *model.Note) error

	// SaveComment 保存评论
	SaveComment(comment *model.Comment) error

	// SaveMedia 保存媒体文件
	SaveMedia(noteID, mediaType, fileName string, content []byte) error

	// GetNoteByID 根据ID获取帖子
	GetNoteByID(noteID string) (*model.Note, error)

	// GetCommentsByNoteID 根据帖子ID获取评论
	GetCommentsByNoteID(noteID string) ([]*model.Comment, error)

	// Close 关闭存储
	Close() error
}

// NewStore 创建存储实例
func NewStore(storeType string, config *config.Config) (Store, error) {
	return NewJSONStore(config) // 默认使用JSON存储
}
