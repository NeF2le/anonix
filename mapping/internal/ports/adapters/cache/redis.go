package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NeF2le/anonix/mapping/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type mappingCache struct {
	ID            uuid.UUID     `json:"id"`
	DekWrapped    []byte        `json:"dek_wrapped,omitempty"`
	CipherText    []byte        `json:"cipher_text,omitempty"`
	TokenTtl      time.Duration `json:"token_ttl,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	Deterministic bool          `json:"deterministic"`
	Reversible    bool          `json:"reversible"`
}

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) SaveMapping(ctx context.Context, mapping *domain.Mapping, ttl time.Duration) error {
	key := fmt.Sprintf("mapping:id:%v", mapping.ID)

	cacheObj := &mappingCache{
		ID:            mapping.ID,
		DekWrapped:    mapping.DekWrapped,
		CipherText:    mapping.CipherText,
		TokenTtl:      mapping.TokenTtl,
		CreatedAt:     mapping.CreatedAt,
		Deterministic: mapping.Deterministic,
		Reversible:    mapping.Reversible,
	}
	payload, err := json.Marshal(cacheObj)
	if err != nil {
		return fmt.Errorf("SaveMapping: failed to marshal mapping: %v", err)
	}

	err = r.client.Set(ctx, key, string(payload), ttl).Err()
	if err != nil {
		return fmt.Errorf("SaveMapping: failed to save mapping in cache: %v", err)
	}
	return nil
}

func (r *RedisAdapter) GetMappingById(ctx context.Context, id uuid.UUID) (*domain.Mapping, error) {
	result := r.client.Get(ctx, fmt.Sprintf("mapping:id:%s", id))
	if err := result.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetMappingById: failed to get mapping: %v", err)
	}

	var mapping domain.Mapping
	if err := json.Unmarshal([]byte(result.Val()), &mapping); err != nil {
		return nil, fmt.Errorf("GetMappingById: failed to unmarshal mapping: %v", err)
	}

	return &mapping, nil
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
