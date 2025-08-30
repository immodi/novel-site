package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/templates/components"
	"immodi/novel-site/internal/http/templates/index"
	"immodi/novel-site/internal/http/templates/novels"
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
	dbNovels, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovels(ctx)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	novels := make([]novels.Novel, len(dbNovels))

	for i, dbNovel := range dbNovels {
		novels[i] = *dbNovelToHomeNovelMapper(dbNovel)
	}

	GenericServiceHandler(w, r, IndexMetaData, index.Index(novels))
}

func dbNovelToHomeNovelMapper(dbNovel repositories.Novel) *novels.Novel {
	return &novels.Novel{
		Name:       dbNovel.Title,
		CoverImage: dbNovel.CoverImage,
	}
}
