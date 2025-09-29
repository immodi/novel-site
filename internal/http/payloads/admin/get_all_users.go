package admin

type AdminGetAllUsersRequest struct {
	Token string `json:"token"`
}

type AdminGetAllUsersResponse struct {
	Users []User `json:"users"`
	Error string `json:"error"`
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	Image     string `json:"image"`
}
