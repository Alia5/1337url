package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Name string
	Pass string
}

var allowedUsers = []User{}

var hmacSampleSecret []byte

func init() {
	hmacSampleSecret = []byte(randomString(256))
}
func AddUsers(users ...User) {
	allowedUsers = append(allowedUsers, users...)
}

func CreateJwt(user string, pass string) (string, error) {
	d1H, err := time.ParseDuration("1h")
	if err != nil {
		return "", err
	}
	d1W, err := time.ParseDuration(fmt.Sprintf("%dh", 7*24))
	if err != nil {
		return "", err
	}

	userValid := false
	for _, u := range allowedUsers {
		if u.Name == user && u.Pass == pass {
			userValid = true
			break
		}
	}
	if !userValid {
		return "", fmt.Errorf("user %s not allowed", user)
	}

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d1W)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-d1H)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(hmacSampleSecret)
}

func ValidateJwt(token string) (bool, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(tok *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok && parsedToken.Valid {
		return ok, nil
	} else {
		return ok, err
	}
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
