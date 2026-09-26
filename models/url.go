package models

import "time"

type URLRecord struct {
	ID             string
	OriginalURL    string
	CreatedAt      time.Time
	AccessedCount  int
	LastAccessedAt time.Time
}
