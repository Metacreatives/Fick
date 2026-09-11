package works

import (
	"encoding/json"
	"net/http"
)

func GetWork(
	repository *Repository,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id := r.PathValue("id")

		work, found, err := findPublicWorkByID(
			r.Context(),
			repository,
			id,
		)

		if err != nil {
			http.Error(
				w,
				"failed to load work",
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
					"error": "work_not_found",
				},
			)

			return
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		if err := json.NewEncoder(w).Encode(work); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)
		}
	}
}
