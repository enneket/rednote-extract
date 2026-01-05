package model

// Note 小红书帖子模型
type Note struct {
	NoteID        string    `json:"note_id" gorm:"primaryKey"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CoverImage    string    `json:"cover_image"`
	Images        []Image   `json:"images"`
	Videos        []Video   `json:"videos"`
	AuthorID      string    `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	AuthorAvatar  string    `json:"author_avatar"`
	Likes         int       `json:"likes"`
	Comments      int       `json:"comments"`
	Collects      int       `json:"collects"`
	Shares        int       `json:"shares"`
	PublishTime   int64     `json:"publish_time"`
	UpdateTime    int64     `json:"update_time"`
	Tags          []string  `json:"tags"`
	Topics        []string  `json:"topics"`
	Location      string    `json:"location"`
	NoteType      string    `json:"note_type"`
	XsecToken     string    `json:"xsec_token"`
	XsecSource    string    `json:"xsec_source"`
	CommentList   []Comment `json:"comment_list"`
	SourceKeyword string    `json:"source_keyword"`
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
}

// Image 帖子图片模型
type Image struct {
	URLDefault   string `json:"url_default"`
	URLThumbnail string `json:"url_thumbnail"`
	URLFull      string `json:"url_full"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	URL          string `json:"url"`
}

// Video 帖子视频模型
type Video struct {
	URL      string `json:"url"`
	CoverURL string `json:"cover_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
}

// NoteURLInfo 帖子URL解析信息
type NoteURLInfo struct {
	NoteID     string `json:"note_id"`
	XsecToken  string `json:"xsec_token"`
	XsecSource string `json:"xsec_source"`
}
