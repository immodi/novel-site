package services

import "log"

const (
	DB_PATH = "db.db"
)

// Services groups all services together
type Services struct {
	DB *DBService
}

func RegisterServices() *Services {
	dbService, err := NewDBService(DB_PATH)
	if err != nil {
		log.Fatal(err)

	}

	return &Services{
		DB: dbService,
	}
}
