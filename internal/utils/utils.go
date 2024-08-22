package utils

import (
	"antrein/bc-queue/model/entity"
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractProjectID(url string) (string, error) {
	re := regexp.MustCompile(`https?://([^.]+)\.antrein\.com`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("URL not registered")
	}
	return matches[1], nil
}

func GenerateJWTToken(key string, claims entity.JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func GenerateSecureRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := crypto.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
