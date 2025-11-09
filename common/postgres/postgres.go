package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     uint16 `yaml:"port" env:"PORT" env-default:"5432"`
	Username string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	Database string `yaml:"db" env:"DB" env-default:"postgres"`
	MaxConns int    `yaml:"max_conns" env:"MAX_CONNS" env-default:"10"`
	MinConns int    `yaml:"min_conns" env:"MIN_CONNS" env-default:"5"`
}

func (c *Config) GetConnString() string {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	return connString
}

func NewPostgresClient(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	connString := config.GetConnString()
	connString += fmt.Sprintf("&pool_max_conns=%d&pool_min_conns=%d",
		config.MaxConns,
		config.MinConns,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func Migrate(ctx context.Context, config *Config, migrationsPath string) error {
	connString := config.GetConnString()

	m, err := migrate.New(migrationsPath, connString)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "migrated successfully")
	return nil
}
