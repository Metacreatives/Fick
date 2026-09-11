package api

import (
	"net/http"
)

func APIRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"/works/",
		http.StripPrefix("/works", WorkRoutes()),
	)

	return mux
}
