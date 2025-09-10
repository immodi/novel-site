package handlers

import (
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"immodi/novel-site/internal/http/templates/errors"
	"net/http"
)

func ServerErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	GenericHandler(w, r, &indexdtostructs.MetaDataStruct{}, errors.InternalServerError())
}

func NotFoundHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	GenericHandler(w, r, &indexdtostructs.MetaDataStruct{}, errors.NotFound())
}
