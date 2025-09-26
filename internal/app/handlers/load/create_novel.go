package load

import (
	"database/sql"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads/load"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strings"
)

func (h *LoadHandler) LoadNovel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, load.LoadNovelResponse{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	// Parse multipart form (limit 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, load.LoadNovelResponse{
			Success: false,
			Message: "failed to parse form: " + err.Error(),
		})
		return
	}

	// Build request object from form fields
	req := load.LoadNovelRequest{
		Title:       r.FormValue("title"),
		Author:      r.FormValue("author"),
		Status:      r.FormValue("status"),
		Description: r.FormValue("description"),
		Genres:      r.Form["genres"],
		Tags:        r.Form["tags"],
	}

	// Handle uploaded cover
	var coverURL string
	file, header, err := r.FormFile("cover_image")
	if err == nil {
		defer file.Close()
		coverURL, err = pkg.SaveUploadedFile(file, header, "static/media")
		println(coverURL)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, load.LoadNovelResponse{
				Success: false,
				Message: "failed to save cover: " + err.Error(),
			})
			return
		}
	}
	req.CoverImage = coverURL

	// Check if novel already exists
	novel, err := h.loadService.GetNovelByExactName(req.Title)
	if err == nil {
		WriteJSON(w, http.StatusConflict, load.LoadNovelResponse{
			Success: false,
			Message: "novel with this title already exists",
			NovelID: novel.ID,
		})
		return
	} else if err != sql.ErrNoRows {
		WriteJSON(w, http.StatusInternalServerError, load.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	// Completed flag
	isCompleted := int64(0)
	if strings.ToLower(req.Status) == "completed" {
		isCompleted = 1
	}

	// Insert into DB
	novel, err = h.loadService.CreateNovel(repositories.CreateNovelParams{
		Title:       req.Title,
		Slug:        pkg.TitleToSlug(req.Title),
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Author:      req.Author,
		AuthorSlug:  pkg.TitleToSlug(req.Author),
		Publisher:   req.Author,
		IsCompleted: isCompleted,
		UpdateTime:  pkg.GetCurrentTimeRFC3339(),
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, load.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	// Add genres in bulk
	if err := h.loadService.AddBulkGenresToNovel(novel.ID, req.Genres); err != nil {
		WriteJSON(w, http.StatusInternalServerError, load.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		go h.loadService.DeleteNovel(novel.ID)
		log.Println(err.Error())
		return
	}

	// Add tags in bulk
	if err := h.loadService.AddBulkTagsToNovel(novel.ID, req.Tags); err != nil {
		WriteJSON(w, http.StatusInternalServerError, load.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		go h.loadService.DeleteNovel(novel.ID)
		log.Println(err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, load.LoadNovelResponse{
		Success: true,
		Message: "novel loaded successfully",
		NovelID: novel.ID,
	})
}
