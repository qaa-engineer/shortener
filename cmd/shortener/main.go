package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qaa-engineer/shortener/internal/handlers"

	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	urlShortenerHandler := handlers.NewURLShortenerHandler()

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
