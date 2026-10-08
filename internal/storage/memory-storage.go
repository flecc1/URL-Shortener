package storage

import (
	"sync"
	"time"
	"url-shortener/models"
)

type MemoryStorage struct {
	data map[string]*models.URLRecord
	mu   sync.RWMutex
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
	m.mu.Lock()
	defer m.mu.Unlock()
	id, err := m.generateUniqID()
	if err != nil {
		return nil, err
	}

	newURLRecord := &models.URLRecord{
		ID:            id,
		CreatedAt:     time.Now(),
		OriginalURL:   url,
		AccessedCount: 0,
	}

	m.data[id] = newURLRecord
	tempURLRecord := *newURLRecord
	return &tempURLRecord, nil
}

func (m *MemoryStorage) Get(id string) (*models.URLRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if record, ok := m.data[id]; ok {
		tempURLRecord := *record
		return &tempURLRecord, nil
	}
	return nil, ErrRecordNotFound
}

func (m *MemoryStorage) DeleteById(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[id]; !ok {
		return ErrRecordNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *MemoryStorage) UpdateById(id, url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if record, ok := m.data[id]; ok {
		record.OriginalURL = url
		return nil
	}
	return ErrRecordNotFound
}

func (m *MemoryStorage) IncrementAccess(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if record, ok := m.data[id]; ok {
		record.AccessedCount++
		record.LastAccessedAt = time.Now()
		return nil
	}
	return ErrRecordNotFound
}

func (m *MemoryStorage) GetAll() ([]*models.URLRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	records := make([]*models.URLRecord, 0, len(m.data))
	for _, record := range m.data {
		tempURLRecord := *record
		records = append(records, &tempURLRecord)
	}
	return records, nil
}
