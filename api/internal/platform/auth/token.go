package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

const Unambiguous = "abcdefghijkmnpqrstuvwxyz23456789"

func RandomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("auth: generate token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func RandomUnambiguous(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("auth: random length must be positive")
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("auth: random unambiguous: %w", err)
	}
	var b strings.Builder
	for _, v := range buf {
		b.WriteByte(Unambiguous[int(v)%len(Unambiguous)])
	}
	return b.String(), nil
}
