package model

// Note 小红书帖子模型
type Note struct {
	NoteID        string     `json:"note_id" gorm:"primaryKey"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	AuthorID      string     `json:"author_id"`
	AuthorName    string     `json:"author_name"`
	LikeCount     int        `json:"like_count"`
	CommentCount  int        `json:"comment_count"`
	CollectCount  int        `json:"collect_count"`
	ShareCount    int        `json:"share_count"`
	PublishTime   int64      `json:"publish_time"`
	UpdateTime    int64      `json:"update_time"`
	Tags          []string   `json:"tags"`
	Topics        []string   `json:"topics"`
	Location      string     `json:"location"`
	NoteType      string     `json:"note_type"`
	XsecToken     string     `json:"xsec_token"`
	XsecSource    string     `json:"xsec_source"`
	Comments      []*Comment `json:"comments"`
	SourceKeyword string     `json:"source_keyword"`
	SourceType    string     `json:"source_type"`
	Keyword       string     `json:"keyword"`
	CreatedAt     int64      `json:"created_at"`
	UpdatedAt     int64      `json:"updated_at"`
}

// NoteURLInfo 帖子URL解析信息
type NoteURLInfo struct {
	NoteID     string `json:"note_id"`
	XsecToken  string `json:"xsec_token"`
	XsecSource string `json:"xsec_source"`
}
