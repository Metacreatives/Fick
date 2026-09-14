package users

import (
	"encoding/json"
	"errors"
	responses "fick/backend/internal/api/generated"
	"io"
	"mime"
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

func RegisterUser(
	repository *Repository,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		mediaType, _, err := mime.ParseMediaType(
			r.Header.Get("Content-Type"),
		)

		if err != nil ||
			mediaType != "application/json" {
			http.Error(
				w,
				"unsupported media type",
				http.StatusUnsupportedMediaType,
			)
			return
		}

		r.Body = http.MaxBytesReader(
			w,
			r.Body,
			4096,
		)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var request responses.RegisterUserRequest

		if err := decoder.Decode(&request); err != nil {
			var sizeError *http.MaxBytesError

			if errors.As(err, &sizeError) {
				http.Error(
					w,
					"request body too large",
					http.StatusRequestEntityTooLarge,
				)
				return
			}

			http.Error(
				w,
				"invalid registration request",
				http.StatusBadRequest,
			)
			return
		}

		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			http.Error(
				w,
				"invalid registration request",
				http.StatusBadRequest,
			)
			return
		}

		if request.Password == nil {
			http.Error(
				w,
				"invalid registration request",
				http.StatusBadRequest,
			)
			return
		}

		user, err := registerUser(
			r.Context(),
			repository,
			request.Username,
			*request.Password,
		)

		if errors.Is(err, ErrInvalidRegistration) ||
			errors.Is(err, ErrUsernameUnavailable) {
			http.Error(
				w,
				"registration failed",
				http.StatusBadRequest,
			)
			return
		}

		if err != nil {
			http.Error(
				w,
				"registration failed",
				http.StatusInternalServerError,
			)
			return
		}

		response := responses.UserResponse{
			Id:        strconv.FormatInt(user.ID, 10),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(response)
	}
}
