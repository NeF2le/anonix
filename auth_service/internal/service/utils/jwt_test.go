package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		ttl    time.Duration
	}{
		{
			name:   "access token",
			userID: "user123",
			ttl:    time.Hour,
		},
	}

	secret := "test-secret"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := GenerateJWT(tt.userID, tt.ttl, secret, false)
			if err != nil {
				t.Fatalf("GenerateJWT() error = %v", err)
			}
			if tokenStr == "" {
				t.Fatalf("GenerateJWT() returned empty token")
			}
			parts := strings.Split(tokenStr, ".")
			if len(parts) != 3 {
				t.Fatalf("GenerateJWT() returned invalid JWT format: %s", tokenStr)
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	secret := "test-secret"
	validTTL := time.Hour
	accessToken, err := GenerateJWT("user123", validTTL, secret, false)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}
	refreshToken, err := GenerateJWT("user123", validTTL, secret, true)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	tests := []struct {
		name        string
		token       string
		secret      string
		wantUser    string
		wantRefresh bool
		wantErr     bool
	}{
		{
			name:    "expired token",
			token:   func() string { tok, _ := GenerateJWT("user123", -1*time.Hour, secret, false); return tok }(),
			secret:  secret,
			wantErr: true,
		},
		{
			name:        "valid access token",
			token:       accessToken,
			secret:      secret,
			wantUser:    "user123",
			wantRefresh: false,
			wantErr:     false,
		},
		{
			name:        "valid refresh token",
			token:       refreshToken,
			secret:      secret,
			wantUser:    "user123",
			wantRefresh: true,
			wantErr:     false,
		},
		{
			name:    "wrong secret",
			token:   accessToken,
			secret:  "not-secret",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "not.a.token",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, isRefresh, _, err := ParseJWT(tt.token, tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
				assert.Equal(t, tt.wantRefresh, isRefresh)
			}
			if err != nil {
				return
			}
		})
	}
}
