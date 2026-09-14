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

	UserRoutes(mux, userRepository)
	WorkRoutes(mux, workRepository, chapterRepository)

	return mux
}
