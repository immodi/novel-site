package homeservice

import (
	"fmt"
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/templates/components"
	"immodi/novel-site/internal/http/templates/index"
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, IndexMetaData, index.Index())
}
