package api

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/works"
	"net/http"
)

func WorkRoutes(
	mux *http.ServeMux,
	workRepository *works.Repository,
	chapterRepository *chapters.Repository,
) http.Handler {
	mux.HandleFunc("GET /works/{id}", works.GetWork(workRepository))
	mux.HandleFunc("GET /works/{work_id}/chapters/{chapter_number}", chapters.GetChapter(chapterRepository))

	return mux
}
