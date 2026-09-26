package storage

import (
	"crypto/rand"
)

const (
	idLength  = 6
	idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateID() (string, error) {
	b := make([]byte, idLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = idCharset[int(b[i])%len(idCharset)]
	}

	return string(b), nil
}
