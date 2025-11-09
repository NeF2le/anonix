package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) DeleteMappingById(ctx context.Context, id uuid.UUID) error {
	result := r.client.Del(ctx, fmt.Sprintf("mapping:id:%s", id))
	if err := result.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("DeleteMappingById: failed to delete mapping: %v", err)
	}
	return nil
}
