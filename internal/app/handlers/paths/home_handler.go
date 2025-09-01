package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/templates/components"
	"immodi/novel-site/internal/http/templates/index"
	homenovelsdto "immodi/novel-site/internal/http/templates/index/components"
	"net/http"
)

var IndexMetaData = &components.MetaDataStruct{
	IsRendering: true,
	Title:       fmt.Sprintf("Read Free Light Novel Online - %s", components.SITE_NAME),
	Description: "We are offering thousands of free books online read! Read novel updated daily: light novel translations, web novel, chinese novel, japanese novel, korean novel, english novel and other novels online.",
	Keywords:    "freewebnovel, novellive, novelfull, mtlnovel, novelupdates, webnovel, korean novel, cultivation novel",
	OgURL:       components.DOMAIN,
	Canonical:   components.DOMAIN,
	CoverImage:  fmt.Sprintf("%s/img/cover.jpg", components.DOMAIN),
	Author:      components.SITE_NAME,
}

type HomeHandler struct {
	dbService *services.DBService
}

func NewHomeHandler(dbService *services.DBService) *HomeHandler {
	return &HomeHandler{dbService: dbService}
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	dbNewestNovels, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNewestHomeNovels(ctx)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbHotNovels, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListHotNovels(ctx)
	})

	newestNovels, err := dbNovelToHomeNovelMapper(dbNewestNovels, h)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hotNovels, err := dbNovelToHomeNovelMapper(dbHotNovels, h)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	GenericServiceHandler(w, r, IndexMetaData, index.Index(hotNovels, newestNovels, newestNovels))
}

func dbNovelToHomeNovelMapper(dbNovels []repositories.Novel, h *HomeHandler) ([]homenovelsdto.HomeNovelDto, error) {
	novels := make([]homenovelsdto.HomeNovelDto, 0, len(dbNovels))

	for _, dbNovel := range dbNovels {

		var dbLatestChapter repositories.Chapter

		dbLatestChapter, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
			return q.GetLatestChapterByNovel(ctx, dbNovel.ID)
		})

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				dbLatestChapter = repositories.Chapter{
					Title: "Chapter doesn't exist",
				}
			} else {
				return nil, err
			}
		}

		dbNovelGenres, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
			return q.ListGenresByNovel(ctx, dbNovel.ID)
		})

		if err != nil {
			return nil, err
		}

		dbChaptersCount, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (int64, error) {
			return q.CountChaptersByNovel(ctx, dbNovel.ID)
		})

		if err != nil {
			return nil, err
		}

		novels = append(novels, homenovelsdto.HomeNovelDto{
			Name:                 dbNovel.Title,
			CoverImage:           dbNovel.CoverImage,
			LastestChapterNumber: int(dbLatestChapter.ChapterNumber),
			LastestChapterName:   dbLatestChapter.Title,
			Status:               dbNovel.Status,
			Genres:               dbNovelGenres,
			LastUpdated:          dbNovel.UpdateTime,
			ChaptersCount:        int(dbChaptersCount),
		})
	}

	return novels, nil
}
