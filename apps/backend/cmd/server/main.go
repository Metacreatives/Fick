package main

import (
	"fick/backend/internal/api"
	"fick/backend/internal/frontend"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	executableDirectory := filepath.Dir(executablePath)

	rendererPath := filepath.Join(
		executableDirectory,
		"..",
		"..",
		"frontend",
		"renderer.mjs",
	)

	rendererProcess, err := frontend.Start(
		rendererPath,
		"http://127.0.0.1:3000",
	)

	if err != nil {
		log.Printf("renderer unavailable: %v", err)
		log.Printf("continuing without frontend rendering")
	} else {
		defer rendererProcess.Close()
	}

	mux := http.NewServeMux()

	mux.Handle(
		"/api/",
		http.StripPrefix("/api", api.APIRoutes()),
	)

	frontendDirectory := filepath.Join(
		executableDirectory,
		"..",
		"..",
		"frontend",
	)

	clientDirectory := filepath.Join(
		frontendDirectory,
		"dist",
		"client",
	)

	rendererBundlePath := filepath.Join(
		frontendDirectory,
		"dist",
		"server",
		"entry-server.js",
	)

	artifactDirectory := filepath.Join(
		executableDirectory,
		"artifacts",
	)

	var frontendHandler http.Handler

	if rendererProcess != nil {
		frontendHandler, err = frontend.FrontendHandler(
			rendererProcess,
			clientDirectory,
			rendererBundlePath,
			artifactDirectory,
		)
		if err != nil {
			log.Printf("frontend unavailable: %v", err)
		}
	}

	if frontendHandler == nil {
		frontendHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"Frontend is unavailable. The backend API is still running.",
				http.StatusServiceUnavailable,
			)
		})
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		frontendHandler.ServeHTTP(w, r)
	}))

	server := &http.Server{
		Addr:              "127.0.0.1:3000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Fick backend is now listening on http://127.0.0.1:3000")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
