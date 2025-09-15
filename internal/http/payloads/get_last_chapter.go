package payloads

type GetLastChapterByIdRequest struct {
	NovelID int64 `json:"novel_id"`
}

type GetLastChapterByNameRequest struct {
	Name string `json:"name"`
}

type GetLastChapterResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message,omitempty"`
	NovelID           int64  `json:"novel_id,omitempty"`
	LastChapterNumber int64  `json:"last_chapter_number,omitempty"`
	LastChapterName   string `json:"last_chapter_name,omitempty"`
}
