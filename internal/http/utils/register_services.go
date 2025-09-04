package utils

import (
	"immodi/novel-site/internal/app/services/auth"
	"immodi/novel-site/internal/app/services/chapters"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/app/services/search"
	"log"
)

const (
	DB_PATH = "db.db"
)

type Services struct {
	DB             *db.DBService
	NovelService   novels.NovelService
	ChapterService chapters.ChapterService
	HomeService    index.HomeService
	AuthService    auth.AuthService
	SearchServie   search.SearchService
}

func RegisterServices() *Services {
	dbService, err := db.NewDBService(DB_PATH)
	if err != nil {
		log.Fatal(err)
	}

	homeService := index.NewHomeService(dbService)

	return &Services{
		DB:             dbService,
		HomeService:    homeService,
		NovelService:   novels.New(dbService),
		ChapterService: chapters.New(dbService),
		AuthService:    auth.New(dbService),
		SearchServie:   search.NewSearchService(dbService, homeService),
	}
}

func (s *Services) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
