package load

import (
	"encoding/json"
	"fmt"
	"immodi/novel-site/internal/http/payloads/load"
	"net/http"
)

func (h *LoadHandler) GetLastChapterById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, load.GetLastChapterResponse{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	var req load.GetLastChapterByIdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}
	defer r.Body.Close()

	dbNovel, err := h.loadService.GetNovelById(req.NovelID)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "invalid novel id",
		})
		return
	}

	dbLastChapter, err := h.loadService.GetLastNovelChapter(dbNovel.ID)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "couldn't get last chapter",
		})
		return
	}

	// Respond success
	WriteJSON(w, http.StatusOK, load.GetLastChapterResponse{
		Success:           true,
		Message:           fmt.Sprintf("got the last chapter of novel %s", dbNovel.Title),
		NovelID:           dbNovel.ID,
		LastChapterNumber: dbLastChapter.ChapterNumber,
		LastChapterName:   dbLastChapter.Title,
	})

}

func (h *LoadHandler) GetLastChapterByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, load.GetLastChapterResponse{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	var req load.GetLastChapterByNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}
	defer r.Body.Close()

	dbNovel, err := h.loadService.GetNovelByExactName(req.Name)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "invalid novel name",
		})
		return
	}

	dbLastChapter, err := h.loadService.GetLastNovelChapter(dbNovel.ID)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, load.GetLastChapterResponse{
			Success: false,
			Message: "couldn't get last chapter",
		})
		return
	}

	// Respond success
	WriteJSON(w, http.StatusOK, load.GetLastChapterResponse{
		Success:           true,
		Message:           fmt.Sprintf("got the last chapter of novel %s", dbNovel.Title),
		NovelID:           dbNovel.ID,
		LastChapterNumber: dbLastChapter.ChapterNumber,
		LastChapterName:   dbLastChapter.Title,
	})

}
