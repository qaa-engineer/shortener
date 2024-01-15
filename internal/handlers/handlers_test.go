package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/qaa-engineer/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLShortenerHandler_GetHandler(t *testing.T) {
	shortURL := "KFDaxKze"
	fullURL := "https://practicum.yandex.ru/learn/go-advanced-self-paced/courses/8bca0296-484d-45dc-b9ab-f01e0f44f9f4/sprints/145736/topics/63027ac1-f19b-405d-bad5-49e3bbddf30b/lessons/572d89a8-1713-457a-927a-90c2280757bc/"
	type fields struct {
		fileStoragePath string
		target          string
		url             string
		store           bool
	}
	type want struct {
		code        int
		contentType string
		location    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test",
			fields: fields{
				fileStoragePath: "",
				target:          shortURL,
				url:             fullURL,
				store:           true,
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    fullURL,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlStorage, _ := storage.NewURLStorage(test.fields.fileStoragePath)
			handler := &URLShortenerHandler{
				URLRepository:   urlStorage,
				BaseResponseURL: "http://example.com/",
			}
			if test.fields.store {
				handler.URLRepository.AddURL(test.fields.target, test.fields.url)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/{id}", http.NoBody)
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", test.fields.target)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeContext))
			handler.GetHandler(w, r)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}

func TestURLShortenerHandler_PostHandler(t *testing.T) {
	fullURL := "https://practicum.yandex.ru/learn/go-advanced-self-paced/courses/8bca0296-484d-45dc-b9ab-f01e0f44f9f4/sprints/145736/topics/63027ac1-f19b-405d-bad5-49e3bbddf30b/lessons/572d89a8-1713-457a-927a-90c2280757bc/"
	type fields struct {
		fileStoragePath string
		target          string
		url             string
	}
	type want struct {
		code        int
		body        string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test",
			fields: fields{
				fileStoragePath: "",
				target:          "/",
				url:             fullURL,
			},
			want: want{
				code:        http.StatusCreated,
				body:        "http://example.com/KFDaxKze",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlStorage, _ := storage.NewURLStorage(test.fields.fileStoragePath)
			handler := &URLShortenerHandler{
				URLRepository:   urlStorage,
				BaseResponseURL: "http://example.com/",
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, test.fields.target, bytes.NewReader([]byte(test.fields.url)))
			handler.PostHandler(w, r)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)

			str := string(resBody)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.body, str)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortenerHandler_PostShortenHandler(t *testing.T) {
	type fields struct {
		fileStoragePath string
		target          string
		shortenRequest  ShortenRequest
	}
	type want struct {
		code            int
		shortenResponse ShortenResponse
		contentType     string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test",
			fields: fields{
				fileStoragePath: "",
				target:          "/api/shorten",
				shortenRequest: ShortenRequest{
					URL: "https://practicum.yandex.ru/learn/go-advanced-self-paced/courses/21df54aa-936c-4845-b435-3c17fc76871a/sprints/145737/topics/769897aa-1f73-4c17-9f2f-9bbbbb4a5cb1/lessons/b36bed70-edce-46d9-84f2-a968e6f2355b/",
				},
			},
			want: want{
				code: http.StatusCreated,
				shortenResponse: ShortenResponse{
					Result: "http://example.com/LnJ2tz3B",
				},
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlStorage, _ := storage.NewURLStorage(tt.fields.fileStoragePath)
			handler := &URLShortenerHandler{
				URLRepository:   urlStorage,
				BaseResponseURL: "http://example.com/",
			}

			w := httptest.NewRecorder()

			body, err := json.Marshal(tt.fields.shortenRequest)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, tt.fields.target, bytes.NewReader(body))
			handler.PostShortenHandler(w, r)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if len(resBody) > 0 {
				var shortenResponse ShortenResponse
				err = json.Unmarshal(resBody, &shortenResponse)
				assert.Equal(t, tt.want.shortenResponse, shortenResponse)
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
