package config

import (
	"fmt"
	"github.com/NeF2le/anonix/common/postgres"
	"github.com/NeF2le/anonix/common/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type AuthServiceConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`

	RedisDB int `yaml:"redis_db" env:"REDIS_DB" env-required:"true"`
}

type Config struct {
	AuthService AuthServiceConfig `yaml:"auth_service" env-prefix:"AUTH_SERVICE_"`
	Postgres    postgres.Config   `yaml:"postgres" env-prefix:"POSTGRES_"`
	Redis       redis.Config      `yaml:"redis" env-prefix:"REDIS_"`

	LogLevel          string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
	MigrationsPath    string        `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	JwtSecret         string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	AccessExpiration  time.Duration `yaml:"access_expiration" env:"ACCESS_EXPIRATION" evn-required:"true"`
	RefreshExpiration time.Duration `yaml:"refresh_expiration" env:"REFRESH_EXPIRATION" evn-required:"true"`
}

func New() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("../.env", &config); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return &config, nil
}
