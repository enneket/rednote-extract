package model

// Comment 小红书评论模型
type Comment struct {
	CommentID       string    `json:"comment_id" gorm:"primaryKey"`
	NoteID          string    `json:"note_id"`
	ParentID        string    `json:"parent_id"`
	UserID          string    `json:"user_id"`
	UserName        string    `json:"user_name"`
	UserAvatar      string    `json:"user_avatar"`
	Content         string    `json:"content"`
	Likes           int       `json:"likes"`
	PublishTime     int64     `json:"publish_time"`
	SubComments     []Comment `json:"sub_comments"`
	SubCommentCount int       `json:"sub_comment_count"`
	CreatedAt       int64     `json:"created_at"`
	UpdatedAt       int64     `json:"updated_at"`
}
