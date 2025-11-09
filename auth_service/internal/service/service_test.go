package service

import (
	"context"
	"github.com/NeF2le/anonix/auth_service/internal/mocks"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const (
	jwtSecret         = "secret"
	refreshExpiration = 5 * time.Second
	accessExpiration  = 5 * time.Second
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		pass      string
		mockSetup func(storage *mocks.StorageRepository)
		expectErr error
	}{
		{
			name:  "happy path",
			login: "test@test.com",
			pass:  "test",
			mockSetup: func(storage *mocks.StorageRepository) {
				storage.
					On("RegisterUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
					Once().
					Return("userID", nil)
			},
			expectErr: nil,
		},
		{
			name:  "user already exists",
			login: "test@test.com",
			pass:  "test",
			mockSetup: func(storage *mocks.StorageRepository) {
				storage.
					On("RegisterUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
					Once().
					Return("", errs.ErrUserAlreadyExists)
			},
			expectErr: errs.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := mocks.NewCacheRepository(t)
			storage := mocks.NewStorageRepository(t)
			if tt.mockSetup != nil {
				tt.mockSetup(storage)
			}

			s := &AuthService{
				storage:           storage,
				cache:             cache,
				jwtSecret:         jwtSecret,
				RefreshExpiration: refreshExpiration,
				AccessExpiration:  accessExpiration,
			}

			userID, err := s.Register(context.Background(), tt.login, tt.pass)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userID)
				assert.Equal(t, "userID", userID)
			}

			cache.AssertExpectations(t)
			storage.AssertExpectations(t)
		})
	}
}
