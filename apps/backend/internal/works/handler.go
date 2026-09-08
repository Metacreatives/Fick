package works

import (
	"encoding/json"
	"net/http"
)

func GetWork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	work, success := findPublicWorkByID(id)

	if success == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "work_not_found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(work); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
