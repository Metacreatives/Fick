package api

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/works"
	"net/http"
)

func WorkRoutes(
	workRepository *works.Repository,
	chapterRepository *chapters.Repository,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", works.GetWork(workRepository))
	mux.HandleFunc("GET /{work_id}/chapters/{chapter_number}", chapters.GetChapter(chapterRepository))

	return mux
}
