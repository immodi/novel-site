package profile

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/http/templates/profile"
	"immodi/novel-site/pkg"
	"log"
	"math"
	"net/http"
	"strconv"
)

func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	pageStr := r.URL.Query().Get("page")

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.profileService.CountBookMarkedNovels(user.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.PROFILE_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.PROFILE_PAGE_LIMIT

	dbNovels, err := h.profileService.ListBookMarkedNovelsPaginated(user.ID, offset, pkg.PROFILE_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	novels, err := index.DbNovelToHomeNovelMapper(dbNovels, h.homeService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	var joinDate = ""
	joinDate, _ = pkg.GetDateFromRFCStr(user.CreatedAt)

	totalPages := int(math.Ceil(float64(totalResults) / float64(pkg.PROFILE_PAGE_LIMIT)))

	data := MapToProfile(user.Username, user.Image, joinDate)
	handlers.GenericHandler(w, r, index.BuildHomeMeta(), profile.Profile(data, novels, currentPage, totalPages))
}
