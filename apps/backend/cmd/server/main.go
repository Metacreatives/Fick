package main

import (
	"fick/backend/internal/chapters"
	"fick/backend/internal/works"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/works/{id}", works.GetWork)
	mux.HandleFunc("GET /api/works/{work_id}/chapters/{chapter_number}", chapters.GetChapter)

	server := &http.Server{
		Addr:              "127.0.0.1:3000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Fick backend is now listening on http://127.0.0.1:3000")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
