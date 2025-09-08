package load

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](r *http.Request) (T, error) {
	var req T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var zero T
		return zero, err
	}
	return req, nil
}

func WriteJSON[T any](w http.ResponseWriter, status int, payload T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
