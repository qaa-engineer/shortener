package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qaa-engineer/shortener/internal/hasher"
	"github.com/qaa-engineer/shortener/internal/storage"
)

type URLShortenerHandler struct {
	URLRepository   storage.URLRepository
	BaseResponseURL string
}

func NewURLShortenerHandler(baseResponseURL string) *URLShortenerHandler {
	return &URLShortenerHandler{
		URLRepository:   storage.NewURLStorage(),
		BaseResponseURL: baseResponseURL,
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

	_, err = w.Write([]byte(handler.BaseResponseURL + shortLink))
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

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (handler *URLShortenerHandler) PostShortenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request ShortenRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	if len(request.URL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("empty request url")
		return
	}

	shortLink, err := hasher.GetShortLink(request.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	handler.URLRepository.AddURL(shortLink, request.URL)
	result := handler.BaseResponseURL + shortLink
	response := ShortenResponse{
		Result: result,
	}

	output, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(output)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}
