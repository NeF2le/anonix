package utils

import (
	"errors"
	"fmt"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	WrappedDek    []byte `json:"wrappedDek"`
	Ciphertext    []byte `json:"ciphertext"`
	Deterministic bool   `json:"deterministic"`
	Reversible    bool   `json:"reversible"`
	jwt.RegisteredClaims
}

func CreateJWTToken(secret string, tokenClaims *domain.TokenClaims) (string, error) {
	claims := jwt.MapClaims{
		"wrappedDek":    tokenClaims.WrappedDek,
		"ciphertext":    tokenClaims.Ciphertext,
		"deterministic": tokenClaims.Deterministic,
		"reversible":    tokenClaims.Reversible,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWTToken(tokenString string, secret string) (*domain.TokenClaims, error) {
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errs.ErrInvalidToken
		}
		return nil, err
	}
	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return &domain.TokenClaims{
		WrappedDek:    claims.WrappedDek,
		Ciphertext:    claims.Ciphertext,
		Deterministic: claims.Deterministic,
		Reversible:    claims.Reversible,
	}, nil
}
