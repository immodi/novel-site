package utils

import (
	"immodi/novel-site/internal/app/services/auth"
	"immodi/novel-site/internal/app/services/chapters"
	comments "immodi/novel-site/internal/app/services/comments"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/load"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/app/services/profile"
	"immodi/novel-site/internal/app/services/search"
	"immodi/novel-site/internal/config"
	"log"
)

type Services struct {
	DB             *db.DBService
	NovelService   novels.NovelService
	ChapterService chapters.ChapterService
	HomeService    index.HomeService
	AuthService    auth.AuthService
	SearchServie   search.SearchService
	LoadService    load.LoadService
	ProfileService profile.ProfileService
	CommentService comments.CommentService
}

func RegisterServices() *Services {
	dbService, err := db.NewDBService(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	homeService := index.NewHomeService(dbService)
	novelService := novels.New(dbService)
	chapterService := chapters.New(dbService, novelService)
	profileSerivce := profile.NewProfileService(dbService, novelService)

	return &Services{
		DB:             dbService,
		HomeService:    homeService,
		NovelService:   novelService,
		ChapterService: chapterService,
		AuthService:    auth.New(dbService),
		SearchServie:   search.NewSearchService(dbService, homeService),
		LoadService:    load.NewLoadService(dbService, novelService, chapterService),
		ProfileService: profileSerivce,
		CommentService: comments.NewCommentService(dbService),
	}
}

func (s *Services) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
