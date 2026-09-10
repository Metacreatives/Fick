package api

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/works"
	"net/http"
)

func WorkRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", works.GetWork)
	mux.HandleFunc("GET /{work_id}/chapters/{chapter_number}", chapters.GetChapter)

	return mux
}
