package auth

import (
	"context"
	"errors"
	"fmt"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
	"strings"
)

type authService struct {
	db *db.DBService
}

func New(db *db.DBService) AuthService {
	return &authService{db: db}
}

func (a *authService) LoginUserWithUsername(username string, password string) (repositories.User, error) {
	return db.ExecuteWithResult(a.db, func(ctx context.Context, q *repositories.Queries) (repositories.User, error) {
		user, err := q.GetUserByUsername(ctx, username)
		if err != nil {
			return repositories.User{}, err
		}

		isValid, err := a.AuthenticateForLogin(username, password)

		if !isValid {
			return repositories.User{}, err
		}

		return user, nil
	})
}

func (a *authService) LoginUserWithEmail(email string, password string) (repositories.User, error) {
	return db.ExecuteWithResult(a.db, func(ctx context.Context, q *repositories.Queries) (repositories.User, error) {
		user, err := q.GetUserByEmail(ctx, email)
		if err != nil {
			return repositories.User{}, err
		}

		isValid, err := a.AuthenticateForLogin(user.Username, password)

		if !isValid {
			return repositories.User{}, err
		}

		return user, nil
	})
}

func (a *authService) RegisterUser(params repositories.CreateUserParams) (int64, error) {
	return db.ExecuteWithResult(a.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		isValid, err := a.AuthenticateForRegister(params.Username, params.PasswordHash)
		if !isValid {
			return 0, err
		}

		userID, err := q.CreateUser(ctx, params)
		if err != nil {
			// Check if it's a unique constraint error
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				if strings.Contains(err.Error(), "users.username") {
					return 0, errors.New("username already taken")
				}
				if strings.Contains(err.Error(), "idx_users_email_lower") {
					return 0, errors.New("email already registered")
				}
				return 0, errors.New("could not store user")
			}

			return 0, err
		}

		return userID, nil
	})
}

func (a *authService) AuthenticateForLogin(username, password string) (bool, error) {
	return db.ExecuteWithResult(a.db, func(ctx context.Context, q *repositories.Queries) (bool, error) {
		u, err := q.GetUserByUsername(ctx, username)
		if err != nil {
			return false, fmt.Errorf("user not found")
		}

		if !pkg.CheckPassword(password, u.PasswordHash) {
			return false, fmt.Errorf("invalid password")
		}

		return true, nil
	})
}

func (a *authService) AuthenticateForRegister(username, password string) (bool, error) {
	return db.ExecuteWithResult(a.db, func(ctx context.Context, q *repositories.Queries) (bool, error) {
		_, err := q.GetUserByUsername(ctx, username)
		if err == nil {
			return false, fmt.Errorf("user already exists")
		}

		return true, nil
	})
}
