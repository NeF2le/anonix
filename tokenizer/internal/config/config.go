package config

import (
	"github.com/NeF2le/anonix/common/tls_helpers"
	"github.com/NeF2le/anonix/common/vault_agent"
	"github.com/ilyakaznacheev/cleanenv"
)

type TokenizerConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`
}

type Config struct {
	Tokenizer  TokenizerConfig    `yaml:"tokenizer" env-prefix:"TOKENIZER_"`
	VaultAgent vault_agent.Config `yaml:"vault_agent" env-prefix:"VAULT_AGENT_"`
	TLS        tls_helpers.Config `yaml:"tls"  env-prefix:"TLS_"`

	ConvergentKey string `yaml:"convergent_key" env:"CONVERGENT_KEY" env-required:"true"`
	DEKBitsLength int    `yaml:"dek_bits_length" env:"DEK_BITS_LENGTH" env-required:"true"`
	LogLevel      string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("../.env", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
