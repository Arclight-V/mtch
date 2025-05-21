package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type BaseClaims struct {
	UserId string `json:"sub"`
	Sid    string `json:"sid"`
}

type AccessClaims struct {
	BaseClaims
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	BaseClaims
	Jti string `json:"jti"`
	Typ string `json:"typ"`
	jwt.RegisteredClaims
}
