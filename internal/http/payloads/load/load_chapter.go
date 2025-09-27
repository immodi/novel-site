package load

type LoadChapterRequest struct {
	NovelID int64 `json:"novel_id"`
}

type LoadChapter struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AppendChapter struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type LoadChapterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
