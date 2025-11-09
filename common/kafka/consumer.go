package kafka

import (
	"context"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type ReaderConfig struct {
	Brokers        []string `yaml:"brokers" env:"BROKERS" env-separator:","`
	CommitInterval int      `yaml:"commit_interval_ms" env:"COMMIT_INTERVAL_MS" env-default:"1000"`
	MaxWorkers     int      `yaml:"max_workers" env:"MAX_WORKERS" env-default:"10"`
}

func NewReader(ctx context.Context, cfg *ReaderConfig, topic, groupID string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: time.Duration(cfg.CommitInterval) * time.Millisecond,
	})
	logger.GetLoggerFromCtx(ctx).Info(ctx, "connected to Kafka topic",
		slog.Any("brokers", cfg.Brokers),
		slog.String("topic", topic),
		slog.String("group_id", groupID),
	)
	return r
}
