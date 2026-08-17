package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	hashVersion = "pbkdf2-sha256"
	iterations  = 210000
	keySize     = 32
	saltSize    = 16
	tokenSize   = 32
)

func HashPassword(password string) (string, error) {
	salt, err := randomBytes(saltSize)
	if err != nil {
		return "", err
	}

	key := deriveKey([]byte(password), salt, iterations, keySize)

	return fmt.Sprintf(
		"%s$%d$%s$%s",
		hashVersion,
		iterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func VerifyPassword(password string, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 || parts[0] != hashVersion {
		return false
	}

	iter, err := strconv.Atoi(parts[1])
	if err != nil || iter <= 0 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}

	actual := deriveKey([]byte(password), salt, iter, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func NewSessionToken() (string, string, error) {
	raw, err := randomBytes(tokenSize)
	if err != nil {
		return "", "", err
	}

	token := base64.RawURLEncoding.EncodeToString(raw)
	return token, HashToken(token), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

func randomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func deriveKey(password []byte, salt []byte, iter int, keyLen int) []byte {
	hashLen := sha256.Size
	numBlocks := (keyLen + hashLen - 1) / hashLen
	output := make([]byte, 0, numBlocks*hashLen)

	for block := 1; block <= numBlocks; block++ {
		u := pbkdfBlock(password, salt, iter, block)
		output = append(output, u...)
	}

	return output[:keyLen]
}

func pbkdfBlock(password []byte, salt []byte, iter int, block int) []byte {
	mac := hmac.New(sha256.New, password)
	mac.Write(salt)
	mac.Write([]byte{
		byte(block >> 24),
		byte(block >> 16),
		byte(block >> 8),
		byte(block),
	})

	u := mac.Sum(nil)
	out := make([]byte, len(u))
	copy(out, u)

	for i := 1; i < iter; i++ {
		mac = hmac.New(sha256.New, password)
		mac.Write(u)
		u = mac.Sum(nil)
		for j := range out {
			out[j] ^= u[j]
		}
	}

	return out
}
