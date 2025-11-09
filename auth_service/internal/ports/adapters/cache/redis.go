package cache

import (
	"context"
	"errors"
	"fmt"
	errs "github.com/NeF2le/anonix/common/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func accessToken(token string) string {
	return fmt.Sprintf("access:%s", token)
}

func refreshToken(token string) string {
	return fmt.Sprintf("refresh:%s", token)
}

type AuthRedisAdapter struct {
	client            *redis.Client
	RefreshExpiration time.Duration
	AccessExpiration  time.Duration
}

func NewAuthRedisAdapter(client *redis.Client, refreshExpiration, accessExpiration time.Duration) *AuthRedisAdapter {
	return &AuthRedisAdapter{
		client:            client,
		RefreshExpiration: refreshExpiration,
		AccessExpiration:  accessExpiration,
	}
}

func (a AuthRedisAdapter) SaveToken(ctx context.Context, token, userID string, refresh bool) error {
	switch refresh {
	case true:
		result := a.client.Set(ctx, refreshToken(token), userID, a.RefreshExpiration)
		if result.Err() != nil {
			return fmt.Errorf("failed to save refresh token: %w", result.Err())
		}
	case false:
		result := a.client.Set(ctx, accessToken(token), userID, a.AccessExpiration)
		if result.Err() != nil {
			return fmt.Errorf("failed to save access token: %w", result.Err())
		}
	}
	return nil
}

func (a AuthRedisAdapter) GetToken(ctx context.Context, token string, refresh bool) (string, error) {
	result := &redis.StringCmd{}
	var err error

	switch refresh {
	case true:
		result = a.client.Get(ctx, refreshToken(token))
		if result.Err() != nil {
			err = fmt.Errorf("failed to get refresh token: %w", result.Err())
		}
	case false:
		result = a.client.Get(ctx, accessToken(token))
		if result.Err() != nil {
			err = fmt.Errorf("failed to get access token: %w", result.Err())
		}
	}
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return "", errs.ErrTokenExpired
		}
		return "", err
	}

	return result.Val(), nil
}

func (a AuthRedisAdapter) DeleteToken(ctx context.Context, token string, refresh bool) error {
	switch refresh {
	case true:
		result := a.client.Del(ctx, refreshToken(token))
		if result.Err() != nil {
			return fmt.Errorf("failed to delete refresh token: %w", result.Err())
		}
	case false:
		result := a.client.Del(ctx, accessToken(token))
		if result.Err() != nil {
			return fmt.Errorf("failed to delete access token: %w", result.Err())
		}
	}
	return nil
}
