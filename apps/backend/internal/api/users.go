package api

import (
	"fick/backend/internal/users"
	"net/http"
)

func UserRoutes(
	userRepository *users.Repository,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", users.GetUser(userRepository))

	return mux
}
