package helper

import (
	"crypto/rand"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func UUID() uuid.UUID {
	return uuid.New()
}

func ShortUUID(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)

	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}

	alphabetLen := byte(len(alphabet))

	for i := 0; i < length; i++ {
		bytes[i] = alphabet[randomBytes[i]%alphabetLen]
	}

	return string(bytes)
}

func GenerateUniqueSlug(name string, withSuffix bool) (string, error) {
	baseSlug := CleanSlug(name)

	if baseSlug == "" {
		return "", eris.New("Slug generation failed")
	}

	if !withSuffix {
		return baseSlug, nil
	}

	suffix := ShortUUID(4)

	return fmt.Sprintf("%s-%s", baseSlug, suffix), nil
}

func CleanSlug(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))

	var buf strings.Builder
	buf.Grow(len(s))

	lastWasDash := false

	for _, r := range s {
		switch {
		// Keep alphanumeric ASCII characters (a-z, 0-9)
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			buf.WriteRune(r)
			lastWasDash = false

		// Replace spaces, underscores, dashes, and whitespace with a single hyphen
		case r == ' ' || r == '-' || r == '_' || unicode.IsSpace(r):
			if !lastWasDash && buf.Len() > 0 {
				buf.WriteRune('-')
				lastWasDash = true
			}

		// Drop all other special characters/punctuation
		default:
			continue
		}
	}

	// Trim trailing dash if present
	result := buf.String()
	return strings.TrimSuffix(result, "-")
}
