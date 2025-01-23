// internal/domain/user.go
package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type JWTClaims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}