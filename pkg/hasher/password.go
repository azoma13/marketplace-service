package hasher

import (
	"crypto/sha512"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type SHA512Hasher struct {
	salt string
}

func NewSHA512Hasher(salt string) *SHA512Hasher {
	return &SHA512Hasher{salt: salt}
}

func (s *SHA512Hasher) Hash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt)))
}
