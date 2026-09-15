package api

import (
	"fick/backend/internal/auth"
	"fick/backend/internal/chapters"
	"fick/backend/internal/users"
	"fick/backend/internal/works"
	"net/http"
)

func APIRoutes(
	workRepository *works.Repository,
	chapterRepository *chapters.Repository,
	userRepository *users.Repository,
	authHandler *auth.Handler,
) http.Handler {
	mux := http.NewServeMux()

	UserRoutes(mux, userRepository, authHandler)
	WorkRoutes(mux, workRepository, chapterRepository)

	return mux
}
