package dto

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ID          string `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}
type StatsResponse struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"created_at"`
	LastAccessedAt string `json:"last_access_at"`
	AccessCount    int    `json:"access_count"`
}
