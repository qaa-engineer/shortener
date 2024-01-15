package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qaa-engineer/shortener/internal/config"
	"github.com/qaa-engineer/shortener/internal/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("text/plain"))

	cfg := config.NewConfiguration()
	urlShortenerHandler, err := handlers.NewURLShortenerHandler(cfg.BaseResponseURL, cfg.FileStoragePath)
	if err != nil {
		panic(err)
	}

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)
	r.Post("/shorten", urlShortenerHandler.PostShortenHandler)

	fmt.Printf("Сервер запушен по адресу %v\n", cfg.BaseResponseURL)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
