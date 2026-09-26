package storage

import (
	"time"
	"url-shortener/models"
)

type MemoryStorage struct {
	data map[string]*models.URLRecord
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

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]*models.URLRecord),
	}
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
