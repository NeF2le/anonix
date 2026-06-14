package utils

import (
	"errors"
	"fmt"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userID string, ttl time.Duration, jwtSecret string, isRefresh bool, roles []string, clearanceLevel int) (string, error) {
	claims := jwt.MapClaims{
		"exp":             time.Now().Add(ttl).Unix(),
		"iat":             time.Now().Unix(),
		"sub":             userID,
		"is_refresh":      isRefresh,
		"roles":           roles,
		"clearance_level": clearanceLevel,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseJWT(tokenString string, jwtSecret string) (string, bool, time.Time, []string, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", false, time.Time{}, nil, 0, errs.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", false, time.Time{}, nil, 0, errs.ErrInvalidToken
		}
		return "", false, time.Time{}, nil, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, time.Time{}, nil, 0, fmt.Errorf("invalid token claims")
	}

	userID, err := claims.GetSubject()
	if err != nil || userID == "" {
		return "", false, time.Time{}, nil, 0, fmt.Errorf("missing or invalid 'sub' in token claims")
	}

	isRefresh, ok := claims["is_refresh"].(bool)
	if !ok {
		return "", false, time.Time{}, nil, 0, fmt.Errorf("missing 'is_refresh' in token claims")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", false, time.Time{}, nil, 0, fmt.Errorf("missing 'exp' in token claims")
	}

	var roles []string
	if raw, ok := claims["roles"].([]interface{}); ok {
		for _, r := range raw {
			if name, ok := r.(string); ok {
				roles = append(roles, name)
			}
		}
	}

	clearanceLevel := 1
	if raw, ok := claims["clearance_level"].(float64); ok {
		clearanceLevel = int(raw)
	}

	return userID, isRefresh, exp.Time, roles, clearanceLevel, nil
}
