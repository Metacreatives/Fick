package main

import (
	"context"
	"fick/backend/internal/api"
	"fick/backend/internal/chapters"
	"fick/backend/internal/config"
	"fick/backend/internal/database"
	dbgen "fick/backend/internal/database/generated"
	"fick/backend/internal/frontend"
	"fick/backend/internal/users"
	"fick/backend/internal/works"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

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
		cfg.Address,
	)

	if err != nil {
		log.Printf("renderer unavailable: %v", err)
		log.Printf("continuing without frontend rendering")
	} else {
		defer rendererProcess.Close()
	}

	ctx := context.Background()

	pool, err := database.Open(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	queries := dbgen.New(pool)

	workRepository := works.NewRepository(queries)
	chapterRepository := chapters.NewRepository(queries)
	userRepository := users.NewRepository(queries)

	mux := http.NewServeMux()

	mux.Handle(
		"/api/",
		http.StripPrefix("/api", api.APIRoutes(
			workRepository,
			chapterRepository,
			userRepository,
		)),
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
		frontendHandler, err = frontend.NewHandler(
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
		Addr:              cfg.Address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Fick backend is now listening on ", cfg.Address)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
