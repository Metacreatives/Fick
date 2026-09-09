package chapters

import (
	"encoding/json"
	"net/http"
)

func GetChapter(w http.ResponseWriter, r *http.Request) {
	work_id := r.PathValue("work_id")
	chapter_number := r.PathValue("chapter_number")

	chapter, success := findChapterByWorkAndNumber(work_id, chapter_number)

	if success == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "chapter_not_found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(chapter); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
