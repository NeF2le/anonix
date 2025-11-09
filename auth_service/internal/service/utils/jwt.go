package utils

import (
	"errors"
	"fmt"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userID string, ttl time.Duration, jwtSecret string, isRefresh bool) (string, error) {
	claims := jwt.MapClaims{
		"exp":        time.Now().Add(ttl).Unix(),
		"iat":        time.Now().Unix(),
		"sub":        userID,
		"is_refresh": isRefresh,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseJWT(tokenString string, jwtSecret string) (string, bool, time.Time, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", false, time.Time{}, errs.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", false, time.Time{}, errs.ErrInvalidToken
		}
		return "", false, time.Time{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, time.Time{}, fmt.Errorf("invalid token claims")
	}

	userID, err := claims.GetSubject()
	if err != nil || userID == "" {
		return "", false, time.Time{}, fmt.Errorf("missing or invalid 'sub' in token claims")
	}

	isRefresh, ok := claims["is_refresh"].(bool)
	if !ok {
		return "", false, time.Time{}, fmt.Errorf("missing 'is_refresh' in token claims")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", false, time.Time{}, fmt.Errorf("missing 'exp' in token claims")
	}

	return userID, isRefresh, exp.Time, nil
}
