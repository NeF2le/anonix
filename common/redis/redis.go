package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"redis"`
	Port     uint16 `yaml:"port" env:"PORT" env-default:"6379"`
	Password string `yaml:"password" env:"PASSWORD"`
	Username string `yaml:"user" env:"USER"`

	PoolSize int `yaml:"pool_size" env:"POOL_SIZE" env-default:"10"`
}

func NewRedisClient(ctx context.Context, cfg *Config, dbNum int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		Username: cfg.Username,
		DB:       dbNum,
		PoolSize: cfg.PoolSize,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetRedisConfigInfo(ctx context.Context, client *redis.Client) string {
	pars := []string{
		"maxmemory",
		"maxmemory-policy",
		"maxmemory-samples",
		"lfu-log-factor",
		"lfu-decay-time",
	}
	result := make(map[string]string, len(pars))
	for _, par := range pars {
		values, err := client.ConfigGet(ctx, par).Result()
		if err != nil {
			continue
		}

		if v, ok := values[par]; ok {
			result[par] = v
		}
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
