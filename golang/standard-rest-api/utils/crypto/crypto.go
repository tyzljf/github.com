package crypto

import (
	"io"
	"crypto/rand"
	"log"
	"encoding/hex"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

//GenerateSalt generates a random salt
func GenerateSalt() string {
	saltBytes := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, saltBytes)
	if err != nil {
		log.Fatal(err)
	}
	salt := make([]byte, 32)
	hex.Decode(salt, saltBytes)
	return string(salt)
}

//HashPassword hashes a string
func HashPassword(password, salt string) string {
	hashPasswordBytes, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		log.Fatal("Unable to hash password")
	}
	hashPassword := make([]byte, 64)
	hex.Decode(hashPassword, hashPasswordBytes)
	return string(hashPassword)
}

func GenerateToken() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	str := base64.URLEncoding.EncodeToString(b)
	return str, nil
}
