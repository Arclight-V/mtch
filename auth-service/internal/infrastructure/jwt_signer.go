package infrastructure

import (
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"github.com/golang-jwt/jwt"
)

type JWTSigner struct {
	key []byte
	alg jwt.SigningMethod
}

func NewJWTSigner(secret []byte) *JWTSigner {
	return &JWTSigner{key: secret, alg: jwt.SigningMethodHS256}
}

func (s *JWTSigner) Sign(c domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(s.alg, jwt.MapClaims{
		"subject": c.UserId,
		"exp":     c.Exp.Unix(),
		"role":    c.Role,
		"issuer":  "auth-service",
	})
	return token.SignedString(s.key)
}
