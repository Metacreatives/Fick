package chapters

import (
	"encoding/json"
	responses "fick/backend/internal/api/generated"
	"fmt"
	"net/http"
)

func GetChapter(repository *Repository) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		workID := r.PathValue("work_id")
		chapterNumber := r.PathValue("chapter_number")

		chapter, found, err := findChapterByWorkAndNumber(
			r.Context(),
			repository,
			workID,
			chapterNumber,
		)

		if err != nil {
			http.Error(
				w,
				"failed to load chapter",
				http.StatusInternalServerError,
			)
			return
		}

		if !found {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusNotFound)

			_ = json.NewEncoder(w).Encode(
				map[string]string{
					"error": "chapter_not_found",
				},
			)

			return
		}

		chapterResponse := responses.ChapterResponse{
			ContentFormat: chapter.ContentFormat,
			ContentRaw:    chapter.ContentRaw,
			CreatedAt:     chapter.CreatedAt.Time,
			Id:            fmt.Sprint(chapter.ID, 10),
			Number:        fmt.Sprint(chapter.Number, 10),
			Title:         chapter.Title,
			UpdatedAt:     chapter.UpdatedAt.Time,
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		if err := json.NewEncoder(w).Encode(chapterResponse); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)
		}
	}
}
