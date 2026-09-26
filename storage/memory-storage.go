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

func (m *MemoryStorage) generateUniqID() (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}
	if _, ok := m.data[id]; ok {
		return m.generateUniqID()
	}
	return id, nil
}

func (m *MemoryStorage) Create(url string) (*models.URLRecord, error) {
	var newURLRecord models.URLRecord
	id, err := m.generateUniqID()
	if err != nil {
		return nil, err
	}
	newURLRecord.ID = id
	newURLRecord.OriginalURL = url
	newURLRecord.CreatedAt = time.Now()
	newURLRecord.AccessedCount = 0
	m.data[id] = &newURLRecord
	return &newURLRecord, nil
}
