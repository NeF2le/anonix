package config

import (
	"fmt"
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type UpstreamNamesConfig struct {
	Mapping   string `yaml:"mapping" env:"MAPPING" env-required:"true"`
	Tokenizer string `yaml:"tokenizer" env:"TOKENIZER" env-required:"true"`
	Auth      string `yaml:"auth" env:"AUTH" env-required:"true"`
}

type UpstreamPortsConfig struct {
	Mapping   string `yaml:"mapping" env:"MAPPING" env-required:"true"`
	Tokenizer string `yaml:"tokenizer" env:"TOKENIZER" env-required:"true"`
	Auth      string `yaml:"auth" env:"AUTH" env-required:"true"`
}

type TimeoutsConfig struct {
	Mapping   time.Duration `yaml:"mapping" env:"MAPPING" env-required:"true"`
	Tokenizer time.Duration `yaml:"tokenizer" env:"TOKENIZER" env-required:"true"`
	Auth      time.Duration `yaml:"auth" env:"AUTH" env-required:"true"`
}

type GrpcPoolConfig struct {
	MaxConnections             int           `yaml:"max_connections" env:"MAX_CONNECTIONS" env-default:"10"`
	MinConnections             int           `yaml:"min_connections" env:"MIN_CONNECTIONS" env-default:"1"`
	MaxRetries                 uint          `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	BaseRetryDelayMilliseconds time.Duration `yaml:"base_retry_delay_milliseconds" env:"BASE_RETRY_DELAY_MILLISECONDS" env-default:"200ms"`
}

type Gateway struct {
	Host string `yaml:"host" env:"HOST" env-required:"true"`
}

type Config struct {
	UpstreamNames UpstreamNamesConfig `yaml:"upstream_names" env-prefix:"UPSTREAM_NAMES_"`
	UpstreamPorts UpstreamPortsConfig `yaml:"upstream_ports" env-prefix:"UPSTREAM_PORTS_"`
	GrpcPool      GrpcPoolConfig      `yaml:"grpc_pool" env-prefix:"GRPC_POOL_"`
	Timeouts      TimeoutsConfig      `yaml:"timeouts" env-prefix:"TIMEOUT_"`
	HTTPPort      int                 `yaml:"http_port" env:"HTTP_PORT" env-default:"8080"`
	Gateway       Gateway             `yaml:"gateway" env-prefix:"GATEWAY_"`
	TLS           tls_helpers.Config  `yaml:"tls" env-prefix:"TLS_"`

	LogLevel              string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
	JWTSecret             string `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"secret"`
	AccessTokenCookieTTL  int    `yaml:"access_token_cookie_ttl" env:"ACCESS_TOKEN_COOKIE_TTL" env-default:"3600"`
	RefreshTokenCookieTTL int    `yaml:"refresh_token_cookie_ttl" env:"REFRESH_TOKEN_COOKIE_TTL" env-default:"36000"`
	Mode                  string `yaml:"mode" env:"MODE" env-required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("../.env", &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}
	return cfg, nil
}
