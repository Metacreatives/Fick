package api

import (
	"fick/backend/internal/users"
	"net/http"
)

func UserRoutes(
	mux *http.ServeMux,
	userRepository *users.Repository,
) http.Handler {
	mux.HandleFunc("POST /users", users.RegisterUser(userRepository))
	mux.HandleFunc("GET /users/{id}", users.GetUser(userRepository))

	return mux
}
