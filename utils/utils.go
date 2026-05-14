package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword accepts a string and returns a hashed password.
// It returns an empty string and error if it failed.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// checkPassword accepts a hashed password and plaintext password, compares if they match.
// Returns true on a match.
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func SignedIn(c *gin.Context) bool {

	auth, ok := c.Request.Context().Value("session-id").(bool)
	if !ok {
		return false
	}
	return auth
}
