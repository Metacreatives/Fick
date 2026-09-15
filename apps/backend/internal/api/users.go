package api

import (
	"fick/backend/internal/auth"
	"fick/backend/internal/users"
	"net/http"
)

func UserRoutes(
	mux *http.ServeMux,
	userRepository *users.Repository,
	authHandler *auth.Handler,
) http.Handler {
	mux.HandleFunc("POST /users", authHandler.RegisterUser)
	mux.HandleFunc("GET /users/{id}", users.GetUser(userRepository))
	mux.HandleFunc("POST /users/sessions", authHandler.CreateSession)

	return mux
}
