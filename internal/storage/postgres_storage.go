package storage

import (
	"database/sql"
	"url-shortener/models"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Create(utl string) (*models.URLRecord, error) {
	query := `
			INSERT INTO urls (original_url) 
			VALUES ($1)
			RETURNING id, original_url, created_at, accessed_count
`
	var newURLRecord models.URLRecord

	err := s.db.QueryRow(query, utl).Scan(
		&newURLRecord.ID,
		&newURLRecord.OriginalURL,
		&newURLRecord.CreatedAt,
		&newURLRecord.AccessedCount)
	if err != nil {
		return nil, err
	}
	return &newURLRecord, nil
}
