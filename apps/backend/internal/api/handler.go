package api

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/users"
	"fick/backend/internal/works"
	"net/http"
)

func APIRoutes(
	workRepository *works.Repository,
	chapterRepository *chapters.Repository,
	userRepository *users.Repository,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"/works/",
		http.StripPrefix("/works", WorkRoutes(workRepository, chapterRepository)),
	)

	mux.Handle(
		"/users/",
		http.StripPrefix("/works", UserRoutes(userRepository)),
	)

	return mux
}
