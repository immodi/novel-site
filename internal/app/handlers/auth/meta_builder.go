package auth

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
)

func BuildAuthMeta(title string) *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		Title:       fmt.Sprintf("%s - %s", title, indexdtostructs.SITE_NAME),
		IsRendering: false,
	}
}
