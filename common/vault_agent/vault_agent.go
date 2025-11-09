package vault_agent

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"strconv"
	"time"
)

type Config struct {
	Host    string        `yaml:"host" env:"HOST" env-required:"true"`
	Port    int           `yaml:"port" env:"PORT" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-required:"true"`
}

func (c *Config) GetAddress() string {
	return "http://" + c.Host + ":" + strconv.Itoa(c.Port)
}

func NewVaultAgent(ctx context.Context, config *Config) (*vault.Client, error) {
	defaultConfig := vault.DefaultConfig()
	defaultConfig.Timeout = config.Timeout
	defaultConfig.Address = config.GetAddress()

	client, err := vault.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}

	if _, err = client.Sys().Health(); err != nil {
		return nil, fmt.Errorf("agent proxy health failed: %w", err)
	}

	return client, err
}
