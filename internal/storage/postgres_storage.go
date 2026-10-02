package storage

import (
	"database/sql"
	"errors"
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

func (s *PostgresStorage) Get(id string) (*models.URLRecord, error) {
	query := `
			SELECT id, original_url, created_at, accessed_count, COALESCE(last_accessed_at, '0001-01-01 00:00:00')
			FROM urls
			WHERE id = $1
`
	var newURLRecord models.URLRecord

	err := s.db.QueryRow(query, id).Scan(
		&newURLRecord.ID,
		&newURLRecord.OriginalURL,
		&newURLRecord.CreatedAt,
		&newURLRecord.AccessedCount,
		&newURLRecord.LastAccessedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("url not found")
		}
		return nil, err
	}
	return &newURLRecord, nil
}

func (s *PostgresStorage) UpdateById(id string, url string) error {
	query := `
			UPDATE urls
			SET original_url = $1
			WHERE id = $2
`
	result, err := s.db.Exec(query, id, url)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("url not found")
	}
	return nil
}
