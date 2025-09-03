package auth

import "immodi/novel-site/internal/app/services/db"

type authService struct {
	db *db.DBService
}

func New(db *db.DBService) AuthService {
	return &authService{db: db}
}
