package sessions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const tokenByteLength = 32

type Token struct {
	Raw  string
	Hash []byte
}

func NewToken() (Token, error) {
	randomBytes := make([]byte, tokenByteLength)

	if _, err := rand.Read(randomBytes); err != nil {
		return Token{}, fmt.Errorf(
			"generate session token: %w",
			err,
		)
	}

	raw := base64.RawURLEncoding.EncodeToString(
		randomBytes,
	)

	hash := HashToken(raw)

	return Token{
		Raw:  raw,
		Hash: hash,
	}, nil
}

func HashToken(raw string) []byte {
	sum := sha256.Sum256([]byte(raw))

	hash := make([]byte, len(sum))
	copy(hash, sum[:])

	return hash
}
