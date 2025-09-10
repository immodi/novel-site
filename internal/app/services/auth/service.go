package auth

import "immodi/novel-site/internal/db/repositories"

type AuthService interface {
	LoginUserWithEmail(email, password string) (repositories.User, error)
	RegisterUser(params repositories.CreateUserParams) (int64, error)
}
