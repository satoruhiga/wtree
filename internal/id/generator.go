package id

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate creates an 8-character random hex ID
func Generate() (string, error) {
	bytes := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
