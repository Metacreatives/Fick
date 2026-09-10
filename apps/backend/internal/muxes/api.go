package muxes

import (
	"fick/backend/internal/muxes/api"
	"net/http"
)

func APIRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"/works/",
		http.StripPrefix("/works", api.WorkRoutes()),
	)

	return mux
}
