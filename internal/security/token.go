package security

import (
	"crypto/rand"
	"encoding/base64"
)

const DEFAULT_TOKEN_SIZE int = 96

func NewToken(size int) (string, error) {
	bytes := make([]byte, size)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
