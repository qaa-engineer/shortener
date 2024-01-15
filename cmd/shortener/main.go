package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qaa-engineer/shortener/internal/config"
	"github.com/qaa-engineer/shortener/internal/handlers"

	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("text/plain"))

	cfg := config.NewConfiguration()

	urlShortenerHandler := handlers.NewURLShortenerHandler(cfg.BaseResponseURL)

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)
	r.Post("/shorten", urlShortenerHandler.PostShortenHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
