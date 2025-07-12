package infrastructure

import (
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	accessTTl  = 15 * time.Minute
	refreshTTl = 30 * 24 * time.Hour
)

type JWTSigner struct {
	accessKey  []byte
	refreshKey []byte
	verifyKey  []byte
	alg        jwt.SigningMethod
}

func NewJWTSigner(accessKye, refreshKey, verifyKey []byte) *JWTSigner {
	return &JWTSigner{accessKey: accessKye, refreshKey: refreshKey, verifyKey: verifyKey, alg: jwt.SigningMethodHS256}
}

func (s *JWTSigner) SignAccess(userId, sid string) (string, error) {
	claims := domain.AccessClaims{
		BaseClaims: domain.BaseClaims{UserId: userId, Sid: sid},
		Roles:      []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTl)),
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(s.alg, claims).SignedString(s.accessKey)
}

func (s *JWTSigner) SignRefresh(userId, sid string) (string, string, error) {
	jtiUUID := uuid.New()
	claims := domain.RefreshClaims{
		BaseClaims: domain.BaseClaims{UserId: userId, Sid: sid},
		Jti:        jtiUUID.String(),
		Typ:        "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTl)),
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(s.alg, claims).SignedString(s.refreshKey)
	return token, jtiUUID.String(), err
}

func (s *JWTSigner) SignVerifyToken(userId string, ttl time.Duration) (string, error) {
	claims := domain.VerifyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Purpose: "verify",
	}
	return jwt.NewWithClaims(s.alg, claims).SignedString(s.verifyKey)
}
