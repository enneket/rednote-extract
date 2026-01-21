package models

type NoteInput struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Comments []string `json:"comments"`
}

type GeneratedNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

type AgentState struct {
	Input          []*NoteInput
	AnalyzedInput  []*AnalyzedInput
	DraftNote      *DraftNote
	FinalNote      *GeneratedNote
	Iteration      int
	Error          error
	ThoughtProcess []string
}

type AnalyzedInput struct {
	MainTopic     string `json:"main_topic"`
	CorePoints    string `json:"core_points"`
	AudienceNeeds string `json:"audience_needs"`
	UsefulInfo    string `json:"useful_info"`
	Style         string `json:"style"`
	Sentiment     string `json:"sentiment"`
	Keywords      string `json:"keywords"`
}

type DraftNote struct {
	Title        string        `json:"title"`
	Outline      string        `json:"outline"`
	Content      string        `json:"content"`
	Tags         string        `json:"tags"`
	WordCount    int           `json:"word_count"`
	Plagiarism   string        `json:"plagiarism_check"`
	ReviewResult *ReviewResult `json:"review_result"`
}

type ReviewResult struct {
	WordCountOK   bool     `json:"word_count_ok"`
	FormatOK      bool     `json:"format_ok"`
	OriginalityOK bool     `json:"originality_ok"`
	ToneOK        bool     `json:"tone_ok"`
	Issues        []string `json:"issues"`
	Suggestions   []string `json:"suggestions"`
	Pass          bool     `json:"pass"`
}

func (d *DraftNote) PlagiarismCheck() bool {
	return d.Plagiarism == "高" || d.Plagiarism == "中"
}

func (d *DraftNote) IsValid() bool {
	return d.WordCount >= MinWordCount &&
		d.WordCount <= MaxWordCount &&
		len(d.Tags) >= 3 &&
		d.PlagiarismCheck()
}

const (
	MaxIterations   = 3
	TargetWordCount = 400
	MinWordCount    = 300
	MaxWordCount    = 500
	MaxPlagiarism   = 0.25
)
