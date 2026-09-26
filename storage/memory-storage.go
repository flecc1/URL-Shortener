package storage

import (
	"sync"
	"time"
	"url-shortener/models"
)

type MemoryStorage struct {
	data  map[string]*models.URLRecord
	mutex sync.RWMutex
}

func (m *MemoryStorage) Create(url string) (*models.URLRecord, error) {
	var newURLRecord models.URLRecord
	newURLRecord.ID = id
	newURLRecord.OriginalURL = url
	newURLRecord.CreatedAt = time.Now()
	newURLRecord.AccessedCount = 0
	m.data[id] = &newURLRecord
	return &newURLRecord, nil
}
