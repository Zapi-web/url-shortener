package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func randomKey() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to create a key: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
