package admin

type AdminGetAllNovelsResponse struct {
	Novels []Novel `json:"novels"`
	Error  string  `json:"error"`
}

type Novel struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Views  int64  `json:"views"`
	Author string `json:"author"`
	Status string `json:"status"`
}
