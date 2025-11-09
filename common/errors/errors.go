package errors

import (
	"errors"
)

var (
	ErrMappingNotFound      = errors.New("mapping not found")
	ErrMappingAlreadyExists = errors.New("mapping already exists")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrMappingExpired       = errors.New("mapping expired")
	ErrTokenExpired         = errors.New("token expired")
	ErrPasswordTooShort     = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong      = errors.New("password must not exceed 72 bytes")
	ErrPasswordNonASCII     = errors.New("password must contain only ASCII characters (no Cyrillic or Unicode)")
	ErrPasswordWeak         = errors.New("password must contain at least one letter and one number")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)
