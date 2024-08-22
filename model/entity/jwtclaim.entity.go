package entity

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	SessionID string `json:"session_id"`
	ProjectID string `json:"project_id"`
	jwt.RegisteredClaims
}
