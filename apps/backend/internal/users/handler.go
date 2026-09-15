package users

import (
	"encoding/json"
	responses "fick/backend/internal/api/generated"
	"net/http"
	"strconv"
)

func GetUser(
	repository *Repository,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id := r.PathValue("id")

		user, found, err := findUserByID(
			r.Context(),
			repository,
			id,
		)

		if err != nil {
			http.Error(
				w,
				"failed to load user",
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
					"error": "user_not_found",
				},
			)

			return
		}

		userResponse := responses.UserResponse{
			CreatedAt: user.CreatedAt.Time,
			Id:        strconv.FormatInt(user.ID, 10),
			UpdatedAt: user.UpdatedAt.Time,
			Username:  user.Username,
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		if err := json.NewEncoder(w).Encode(userResponse); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)
		}
	}
}
