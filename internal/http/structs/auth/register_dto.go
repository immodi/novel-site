package authdtostructs

type RegisterDTO struct {
	Username string   // pre-fill username field on error
	Email    string   // pre-fill email field on error
	Errors   []string // error messages to display
}
