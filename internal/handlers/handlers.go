package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/qaa-engineer/shortener/internal/hasher"
	"github.com/qaa-engineer/shortener/internal/storage"

	"io"
	"net/http"
)

type URLShortenerHandler struct {
	URLRepository storage.URLRepository
}

func NewURLShortenerHandler() *URLShortenerHandler {
	return &URLShortenerHandler{
		URLRepository: storage.NewURLStorage(),
	}
}

func (handler *URLShortenerHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	w.Header().Set("Content-Type", "text/plain")

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := string(body)
	shortLink, err := hasher.GetShortLink(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	handler.URLRepository.AddURL(shortLink, res)

	w.WriteHeader(http.StatusCreated)
	host := r.Host

	_, err = w.Write([]byte("http://" + host + "/" + shortLink))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *URLShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, ok := handler.URLRepository.GetURL(id)

	w.Header().Set("Content-Type", "text/plain")

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
