package config

import (
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type MappingCleanerConfig struct {
	RedisDB  int           `yaml:"redis_db" env:"REDIS_DB" env-required:"true"`
	Cooldown time.Duration `yaml:"cooldown" env:"COOLDOWN" env-required:"true"`
}

type Config struct {
	Postgres       postgres.Config      `yaml:"postgres" env-prefix:"POSTGRES_"`
	Redis          redis.Config         `yaml:"redis" env-prefix:"REDIS_"`
	MappingCleaner MappingCleanerConfig `yaml:"mapping_cleaner" env-prefix:"MAPPING_CLEANER_"`

	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("../.env", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
