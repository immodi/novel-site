package admin

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token      string `json:"token"`
	Username   string `json:"username"`
	CoverImage string `json:"coverImage"`
	Error      string `json:"error"`
}
