package api

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/works"
	"net/http"
)

func APIRoutes(
	workRepository *works.Repository,
	chapterRepository *chapters.Repository,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"/works/",
		http.StripPrefix("/works", WorkRoutes(workRepository, chapterRepository)),
	)

	return mux
}
