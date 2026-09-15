package auth

import (
	"encoding/json"
	"errors"
	responses "fick/backend/internal/api/generated"
	"fick/backend/internal/sessions"
	"io"
	"mime"
	"net/http"
	"strconv"
)

const maxRequestBodyBytes int64 = 4 * 1024

var (
	errUnsupportedMediaType = errors.New(
		"unsupported media type",
	)

	errRequestBodyTooLarge = errors.New(
		"request body too large",
	)

	errInvalidRequest = errors.New(
		"invalid request",
	)
)

type Handler struct {
	service *Service
	cookies sessions.CookieConfig
}

func NewHandler(
	service *Service,
	cookies sessions.CookieConfig,
) *Handler {
	return &Handler{
		service: service,
		cookies: cookies,
	}
}

func (h *Handler) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request responses.RegisterUserRequest

	if err := decodeJSON(
		w,
		r,
		&request,
	); err != nil {
		writeRequestError(w, err)
		return
	}

	if request.Password == nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"registration_failed",
		)
		return
	}

	result, err := h.service.Register(
		r.Context(),
		request.Username,
		*request.Password,
	)

	if errors.Is(
		err,
		ErrRegistrationFailed,
	) {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"registration_failed",
		)
		return
	}

	if err != nil {
		writeJSONError(
			w,
			http.StatusInternalServerError,
			"registration_failed",
		)
		return
	}

	h.cookies.Set(
		w,
		result.Token,
	)

	response := responses.UserResponse{
		Id: strconv.FormatInt(
			result.User.ID,
			10,
		),

		Username: result.User.Username,

		CreatedAt: result.User.CreatedAt.Time,

		UpdatedAt: result.User.UpdatedAt.Time,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	_ = json.NewEncoder(w).Encode(
		response,
	)
}

func (h *Handler) CreateSession(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request responses.CreateSessionRequest

	if err := decodeJSON(
		w,
		r,
		&request,
	); err != nil {
		writeRequestError(w, err)
		return
	}

	if request.Password == nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
		return
	}

	token, err := h.service.Login(
		r.Context(),
		request.Username,
		*request.Password,
	)

	if errors.Is(
		err,
		ErrInvalidCredentials,
	) {
		writeJSONError(
			w,
			http.StatusUnauthorized,
			"invalid_credentials",
		)
		return
	}

	if err != nil {
		writeJSONError(
			w,
			http.StatusInternalServerError,
			"authentication_failed",
		)
		return
	}

	h.cookies.Set(
		w,
		token,
	)

	w.WriteHeader(
		http.StatusNoContent,
	)
}

func decodeJSON(
	w http.ResponseWriter,
	r *http.Request,
	target any,
) error {
	mediaType, _, err := mime.ParseMediaType(
		r.Header.Get("Content-Type"),
	)

	if err != nil ||
		mediaType != "application/json" {
		return errUnsupportedMediaType
	}

	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxRequestBodyBytes,
	)

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		var sizeError *http.MaxBytesError

		if errors.As(
			err,
			&sizeError,
		) {
			return errRequestBodyTooLarge
		}

		return errInvalidRequest
	}

	if err := decoder.Decode(
		&struct{}{},
	); err != io.EOF {
		var sizeError *http.MaxBytesError

		if errors.As(
			err,
			&sizeError,
		) {
			return errRequestBodyTooLarge
		}

		return errInvalidRequest
	}

	return nil
}

func writeRequestError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(
		err,
		errUnsupportedMediaType,
	):
		writeJSONError(
			w,
			http.StatusUnsupportedMediaType,
			"unsupported_media_type",
		)

	case errors.Is(
		err,
		errRequestBodyTooLarge,
	):
		writeJSONError(
			w,
			http.StatusRequestEntityTooLarge,
			"request_body_too_large",
		)

	default:
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
	}
}

func writeJSONError(
	w http.ResponseWriter,
	status int,
	code string,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.Header().Set(
		"Cache-Control",
		"no-store",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		responses.ErrorResponse{
			Error: code,
		},
	)
}
