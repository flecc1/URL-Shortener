package storage

import (
	"errors"
	"url-shortener/models"
)

type Storage interface {
	Create(url string) (*models.URLRecord, error)
	Get(id string) (*models.URLRecord, error)
	UpdateById(id, url string) error
	DeleteById(id string) error
	IncrementAccess(id string) error
	GetAll() ([]*models.URLRecord, error)
}

var UrlNotFoundError = errors.New("url not found")
