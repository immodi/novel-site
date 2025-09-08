package load

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads"
	"immodi/novel-site/pkg"
)

func (h *LoadHandler) LoadNovel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, payloads.LoadNovelResponse{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	req, err := DecodeJSON[payloads.LoadNovelRequest](r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "failed to decode request",
		})
		log.Println(err.Error())
		return
	}

	// Check if novel already exists
	novel, err := h.loadService.GetNovelByExactName(req.Title)
	if err == nil {
		WriteJSON(w, http.StatusConflict, payloads.LoadNovelResponse{
			Success: false,
			Message: "novel with this title already exists",
			NovelID: novel.ID,
		})
		return
	} else if err != sql.ErrNoRows {
		// Real DB error
		WriteJSON(w, http.StatusInternalServerError, payloads.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	isCompleted := int64(0)
	if strings.ToLower(req.Status) == "completed" {
		isCompleted = 1
	}

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
		WriteJSON(w, http.StatusInternalServerError, payloads.LoadNovelResponse{
			Success: false,
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	for _, genre := range req.Genres {
		err = h.loadService.AddGenreToNovel(novel.ID, genre)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, payloads.LoadNovelResponse{
				Success: false,
				Message: err.Error(),
			})
			go h.loadService.DeleteNovel(novel.ID)
			log.Println(err.Error())
			return
		}
	}

	for _, tag := range req.Tags {
		err = h.loadService.AddTagToNovel(novel.ID, tag)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, payloads.LoadNovelResponse{
				Success: false,
				Message: err.Error(),
			})
			go h.loadService.DeleteNovel(novel.ID)
			log.Println(err.Error())
			return
		}
	}

	WriteJSON(w, http.StatusCreated, payloads.LoadNovelResponse{
		Success: true,
		Message: "novel loaded successfully",
		NovelID: novel.ID,
	})
}

func (h *LoadHandler) LoadChapter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, payloads.LoadNovelResponse{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	err := r.ParseMultipartForm(20 << 20) // 20 MB max
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "failed to parse form (file too large?)",
		})
		log.Println(err.Error())
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("file") // the field name should be 'file'
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "failed to get file",
		})
		log.Println(err.Error())
		return
	}
	defer file.Close()

	// 2. Get the metadata JSON
	metadataStr := r.FormValue("metadata")
	var req payloads.LoadChapterRequest
	if err := json.Unmarshal([]byte(metadataStr), &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "invalid metadata JSON",
		})
		log.Println(err.Error())
		return
	}

	// Decode JSON from file
	var chapters []payloads.LoadChapter
	if err := json.NewDecoder(file).Decode(&chapters); err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "invalid JSON file",
		})
		log.Println(err.Error())
		return
	}

	var createParams []repositories.CreateChapterParams
	for i, ch := range chapters {
		createParams = append(createParams, repositories.CreateChapterParams{
			NovelID:       req.NovelID,
			ChapterNumber: int64(i + 1),
			Title:         ch.Title,
			Content:       ch.Content,
		})
	}

	if err := h.loadService.CreateBulkChapters(createParams); err != nil {
		WriteJSON(w, http.StatusInternalServerError, payloads.LoadNovelResponse{
			Success: false,
			Message: "failed to insert chapters",
		})
		log.Println(err.Error())
		return
	}

	// Respond success
	WriteJSON(w, http.StatusOK, payloads.LoadNovelResponse{
		Success: true,
		Message: fmt.Sprintf("Loaded %d chapters", len(chapters)),
	})
}
