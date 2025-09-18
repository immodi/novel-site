package load

import (
	"encoding/json"
	"fmt"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
)

func (h *LoadHandler) AppendChapters(w http.ResponseWriter, r *http.Request) {
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
	var chapters []payloads.AppendChapter
	if err := json.NewDecoder(file).Decode(&chapters); err != nil {
		WriteJSON(w, http.StatusBadRequest, payloads.LoadNovelResponse{
			Success: false,
			Message: "invalid JSON file",
		})
		log.Println(err.Error())
		return
	}

	var createParams []repositories.CreateChapterParams
	for _, ch := range chapters {
		createParams = append(createParams, repositories.CreateChapterParams{
			NovelID:       req.NovelID,
			ChapterNumber: int64(ch.Number),
			Title:         ch.Title,
			ReleaseDate:   pkg.GetCurrentTimeRFC3339(),
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
