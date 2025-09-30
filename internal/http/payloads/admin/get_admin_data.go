package admin

type AdminGetAdminDataResponse struct {
	Username   string `json:"username"`
	CoverImage string `json:"coverImage"`
	Error      string `json:"error"`
}
