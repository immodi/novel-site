package admin

type AdminGetAllNovelChaptersResponse struct {
	Chapters []Chapter `json:"chapters"`
	Error    string    `json:"error"`
}

type Chapter struct {
	ID          int64  `json:"id"`
	NovelID     int64  `json:"novelId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ReleaseDate string `json:"releaseDate"`
}
