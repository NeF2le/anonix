package config

import (
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type MappingConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`

	RedisDB  int           `yaml:"redis_db" env:"REDIS_DB" env-required:"true"`
	CacheTtl time.Duration `yaml:"cache_ttl" env:"CACHE_TTL" env-default:"10h"`
}

type Config struct {
	Postgres postgres.Config    `yaml:"postgres" env-prefix:"POSTGRES_"`
	Redis    redis.Config       `yaml:"redis" env-prefix:"REDIS_"`
	Mapping  MappingConfig      `yaml:"mapping" env-prefix:"MAPPING_"`
	TLS      tls_helpers.Config `yaml:"tls" env-prefix:"TLS_"`

	LogLevel       string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("../.env", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
