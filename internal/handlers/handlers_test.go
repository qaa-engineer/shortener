package handlers

import (
	"bytes"
	"context"
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
		URLRepository storage.URLRepository
		target        string
		url           string
		store         bool
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
				URLRepository: storage.NewURLStorage(),
				target:        shortURL,
				url:           fullURL,
				store:         true,
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
			handler := &URLShortenerHandler{
				URLRepository: test.fields.URLRepository,
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
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}

func TestURLShortenerHandler_PostHandler(t *testing.T) {
	fullURL := "https://practicum.yandex.ru/learn/go-advanced-self-paced/courses/8bca0296-484d-45dc-b9ab-f01e0f44f9f4/sprints/145736/topics/63027ac1-f19b-405d-bad5-49e3bbddf30b/lessons/572d89a8-1713-457a-927a-90c2280757bc/"
	type fields struct {
		URLRepository storage.URLRepository
		target        string
		url           string
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
				URLRepository: storage.NewURLStorage(),
				target:        "/",
				url:           fullURL,
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
			handler := &URLShortenerHandler{
				URLRepository: test.fields.URLRepository,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, test.fields.target, bytes.NewReader([]byte(test.fields.url)))
			handler.PostHandler(w, r)

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			str := string(resBody)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.body, str)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
