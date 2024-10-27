package helper

import (
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

type UuidHelper struct{}

func NewUuidHelper() *UuidHelper {
	return &UuidHelper{}
}

func (helper *UuidHelper) Generate() uuid.UUID {
	return uuid.New()
}

func (helper *UuidHelper) GenerateShort() string {
	googleUuid := uuid.New()

	b64 := base64.RawURLEncoding.EncodeToString(googleUuid[:6])

	return strings.TrimRight(b64, "=")
}
