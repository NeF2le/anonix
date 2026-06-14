package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/anonix/auth_service/internal/domain"
	"github.com/NeF2le/anonix/auth_service/internal/ports"
	"github.com/NeF2le/anonix/auth_service/internal/service/utils"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	storage           ports.StorageRepository
	cache             ports.CacheRepository
	jwtSecret         string
	RefreshExpiration time.Duration
	AccessExpiration  time.Duration
}

func NewAuthService(
	storage ports.StorageRepository,
	cache ports.CacheRepository,
	jwtSecret string,
	RefreshExpiration time.Duration,
	AccessExpiration time.Duration,
) *AuthService {
	return &AuthService{
		storage:           storage,
		cache:             cache,
		jwtSecret:         jwtSecret,
		RefreshExpiration: RefreshExpiration,
		AccessExpiration:  AccessExpiration,
	}
}

func (s *AuthService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errs.ErrInvalidCredentials
	}

	if err = s.storage.DeleteUser(ctx, userUUID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) AssignRole(ctx context.Context, userId string, roleId int) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errs.ErrInvalidCredentials
	}

	if err = s.storage.AssignRole(ctx, userUUID, roleId); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RemoveRole(ctx context.Context, userId string, roleId int) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errs.ErrInvalidCredentials
	}

	if err = s.storage.RemoveRole(ctx, userUUID, roleId); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to remove role for user", err,
			slog.String("userId", userId))
		return err
	}

	return nil
}

func (s *AuthService) UpdateClearanceLevel(ctx context.Context, userId string, level int) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errs.ErrInvalidCredentials
	}

	if err = s.storage.UpdateClearanceLevel(ctx, userUUID, level); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetRolesList(ctx context.Context) ([]*domain.Role, error) {
	roles, err := s.storage.GetRolesList(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get roles list", err)
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "got roles list")

	return roles, nil
}

func (s *AuthService) Register(ctx context.Context, login string, pass string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userID, err := s.storage.RegisterUser(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return "", errs.ErrUserAlreadyExists
		}
		return "", err
	}

	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, login string, pass string) (string, string, string, error) {
	user, err := s.storage.LoginUser(ctx, login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return "", "", "", errs.ErrInvalidCredentials
		}
		return "", "", "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		return "", "", "", errs.ErrInvalidCredentials
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return "", "", "", errs.ErrInvalidCredentials
	}

	domainRoles, err := s.storage.GetUserRoles(ctx, userUUID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get user roles: %w", err)
	}
	roleNames := make([]string, 0, len(domainRoles))
	for _, r := range domainRoles {
		roleNames = append(roleNames, r.Name)
	}

	accessToken, err := utils.GenerateJWT(user.ID, s.AccessExpiration, s.jwtSecret, false, roleNames, user.ClearanceLevel)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := utils.GenerateJWT(user.ID, s.RefreshExpiration, s.jwtSecret, true, roleNames, user.ClearanceLevel)
	if err != nil {
		return "", "", "", err
	}

	err = s.cache.SaveToken(ctx, accessToken, user.ID, false)
	if err != nil {
		return "", "", "", err
	}
	err = s.cache.SaveToken(ctx, refreshToken, user.ID, true)
	if err != nil {
		return "", "", "", err
	}

	return user.ID, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	sub, ifRefresh, exp, roleNames, clearanceLevel, err := utils.ParseJWT(refreshToken, s.jwtSecret)
	if err != nil {
		if errors.Is(err, errs.ErrTokenExpired) {
			return "", "", errs.ErrTokenExpired
		}
		if errors.Is(err, errs.ErrInvalidToken) {
			return "", "", errs.ErrInvalidToken
		}
		return "", "", err
	}

	if !ifRefresh {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "token is not refresh token")
		return "", "", errs.ErrInvalidToken
	}

	if time.Now().After(exp) {
		return "", "", errs.ErrTokenExpired
	}

	userID, err := s.cache.GetToken(ctx, refreshToken, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch refresh token: %w", err)
	}

	if userID != sub {
		logger.GetLoggerFromCtx(ctx).Warn(ctx,
			"refresh token user id does not match",
			slog.String("user", userID),
			slog.String("sub", sub),
		)
		return "", "", errs.ErrInvalidToken
	}

	accessToken, err := utils.GenerateJWT(userID, s.AccessExpiration, s.jwtSecret, true, roleNames, clearanceLevel)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateJWT(userID, s.RefreshExpiration, s.jwtSecret, false, roleNames, clearanceLevel)
	if err != nil {
		return "", "", err
	}

	err = s.cache.SaveToken(ctx, accessToken, userID, false)
	if err != nil {
		return "", "", err
	}

	err = s.cache.SaveToken(ctx, refreshToken, userID, true)
	if err != nil {
		return "", "", err
	}

	err = s.cache.DeleteToken(ctx, refreshToken, true)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) GetUserRoles(ctx context.Context, userId string) ([]*domain.Role, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	return s.storage.GetUserRoles(ctx, userUUID)
}

func (s *AuthService) IsAdmin(ctx context.Context, userId string) (bool, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return false, errs.ErrInvalidCredentials
	}

	result, err := s.storage.IsAdminCheck(ctx, userUUID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return false, errs.ErrUserNotFound
		}
		return false, err
	}

	return result, nil
}
