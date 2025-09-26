package admin

func authenticateLoginRequest(email, password string) []string {
	var errors []string

	if email == "" {
		errors = append(errors, "Email is required")
	}
	if password == "" {
		errors = append(errors, "Password is required")
	}

	return errors
}
