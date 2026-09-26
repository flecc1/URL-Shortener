package storage

import (
	"time"
	"url-shortener/models"
)

type MemoryStorage map[string]*models.URLRecord

func (m MemoryStorage) Create(url string) (*models.URLRecord, error) {
	var newURLRecord models.URLRecord
	newURLRecord.OriginalURL = url
	newURLRecord.CreatedAt = time.Now()
	newURLRecord.AccessedCount = 0
	// todo напсиать функцию генерации id(short url)
	m[url] = &newURLRecord
	return &newURLRecord, nil
}
