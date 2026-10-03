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

func (m *MemoryStorage) Get(id string) (*models.URLRecord, error) {
	if record, ok := m.data[id]; ok {
		return record, nil
	}
	return nil, UrlNotFoundError
}

func (m *MemoryStorage) DeleteById(id string) error {
	if _, ok := m.data[id]; !ok {
		return UrlNotFoundError
	}
	delete(m.data, id)
	return nil
}

func (m *MemoryStorage) UpdateById(id, url string) error {
	if record, ok := m.data[id]; ok {
		record.OriginalURL = url
		return nil
	}
	return UrlNotFoundError
}

func (m *MemoryStorage) IncrementAccess(id string) error {
	if record, ok := m.data[id]; ok {
		record.AccessedCount++
		record.LastAccessedAt = time.Now()
		return nil
	}
	return UrlNotFoundError
}

func (m *MemoryStorage) GetAll() ([]*models.URLRecord, error) {
	records := make([]*models.URLRecord, 0, len(m.data))
	for _, record := range m.data {
		records = append(records, record)
	}
	return records, nil
}
