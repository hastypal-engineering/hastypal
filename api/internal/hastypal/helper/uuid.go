package helper

import (
	"crypto/rand"
	"github.com/google/uuid"
)

type UuidHelper struct{}

func NewUuidHelper() *UuidHelper {
	return &UuidHelper{}
}

func (helper *UuidHelper) Generate() uuid.UUID {
	return uuid.New()
}

func (helper *UuidHelper) GenerateShort() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

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
