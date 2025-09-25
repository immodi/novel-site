package filter

import (
	"immodi/novel-site/internal/http/templates/filter"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

func (h *FilterHandler) FilterTags(w http.ResponseWriter, r *http.Request) {
	tag := strings.ToLower(r.URL.Query().Get("tag"))
	h.tagList(tag, w, r)
}

func (h *FilterHandler) FilterExcludedTags(w http.ResponseWriter, r *http.Request) {
	tag := strings.ToLower(r.URL.Query().Get("tagExcluded"))
	h.tagList(tag, w, r)
}

func (h *FilterHandler) tagList(tag string, w http.ResponseWriter, r *http.Request) {
	if len(tag) < 3 {
		cmp := filter.TagSelectorWarning()
		templ.Handler(cmp).ServeHTTP(w, r)
		return
	}

	tags, err := h.novelService.FilterTagsByName(tag)
	if err != nil {
		log.Println(err.Error())
		tags = []string{}
	}

	for _, t := range tags {
		log.Println(t)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cmp := filter.TagSelector(tags)
	templ.Handler(cmp).ServeHTTP(w, r)
}
