package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type URLRepository interface {
	AddURL(shortLink, originalLink string)
	GetURL(shortLink string) (string, bool)
}

type URLStorage struct {
	mu   sync.Mutex
	Urls map[string]string
}

type URLEntity struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLStorage(fileStoragePath string) (*URLStorage, error) {

	urls := make(map[string]string)
	if len(fileStoragePath) != 0 {
		fileStorage, err := createFileStorage(fileStoragePath)
		defer func(fileStorage *os.File) {
			err = fileStorage.Close()
		}(fileStorage)

		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileStorage)
		var urlEntity URLEntity
		for scanner.Scan() {

			bytes := scanner.Bytes()
			if len(bytes) == 0 {
				continue
			}
			err := json.Unmarshal(bytes, &urlEntity)

			if err != nil {
				return nil, err
			}

			urls[urlEntity.ShortURL] = urlEntity.OriginalURL
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return &URLStorage{
		Urls: make(map[string]string),
	}, nil
}

func createFileStorage(p string) (*os.File, error) {
	if _, err := os.Stat(p); err == nil {
		file, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (s *URLStorage) AddURL(shortLink, originalLink string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Urls[shortLink] = originalLink
}

func (s *URLStorage) GetURL(shortLink string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.Urls[shortLink]

	return url, ok
}
