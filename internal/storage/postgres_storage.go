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

func (s *PostgresStorage) Create(url string) (*models.URLRecord, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	query := `
			INSERT INTO urls (id, original_url) 
			VALUES ($1, $2)
			RETURNING id, original_url, created_at, accessed_count
`
	var newURLRecord models.URLRecord

	err = s.db.QueryRow(query, id, url).Scan(
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
			return nil, UrlNotFoundError
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
	result, err := s.db.Exec(query, url, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return UrlNotFoundError
	}
	return nil
}

func (s *PostgresStorage) DeleteById(id string) error {
	query := `DELETE FROM urls WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return UrlNotFoundError
	}
	return nil
}

func (s *PostgresStorage) IncrementAccess(id string) error {
	query := `
			UPDATE urls 
			SET accessed_count = accessed_count + 1, last_accessed_at = NOW()
			WHERE id = $1
`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return UrlNotFoundError
	}
	return nil
}
