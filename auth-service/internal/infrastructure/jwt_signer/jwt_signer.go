package jwt_signer

import (
	"errors"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	accessTTl  = 15 * time.Minute
	refreshTTl = 30 * 24 * time.Hour
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrWrongPurpose = errors.New("wrong token purpose")
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

// SignVerifyToken sign verify token. Return token, jti, error
func (s *JWTSigner) SignVerifyToken(userId string, ttl time.Duration) (domain.VerifyTokenIssue, string, error) {
	jti := uuid.NewString()
	expiresAt := jwt.NewNumericDate(time.Now().Add(ttl))

	claims := domain.VerifyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			ID:        jti,
		},
		Purpose: "verify",
	}

	tok, err := jwt.NewWithClaims(s.alg, claims).SignedString(s.verifyKey)
	verifyTokenIssue := domain.VerifyTokenIssue{JTI: jti, UserID: userId, ExpiresAt: expiresAt.Time}
	return verifyTokenIssue, tok, err
}

func (s *JWTSigner) ParseVerifyToken(tokenStr string) (domain.VerifyEmailToken, error) {
	var claims domain.VerifyClaims

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{s.alg.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(2*time.Second),
	)

	tok, err := parser.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.verifyKey, nil
	})
	if err != nil || !tok.Valid {
		return domain.VerifyEmailToken{}, ErrInvalidToken
	}
	if claims.Purpose != "verify" {
		return domain.VerifyEmailToken{}, ErrWrongPurpose
	}

	return domain.VerifyEmailToken{JTI: claims.ID, UserID: claims.Subject}, nil
}
