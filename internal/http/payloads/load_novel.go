package payloads

type LoadNovelRequest struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Genres      []string `json:"genres"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	CoverImage  string   `json:"cover_image"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
}

type LoadNovelResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	NovelID int64  `json:"novel_id,omitempty"`
}
