package admin

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token  string   `json:"token"`
	Errors []string `json:"errors"`
}
