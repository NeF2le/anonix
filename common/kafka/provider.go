package kafka

import (
	"context"
	"github.com/NeF2le/anonix/common/logger"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type WriterConfig struct {
	Brokers []string `yaml:"brokers" env:"BROKERS" env-separator:","`
}

func NewWriter(ctx context.Context, cfg *WriterConfig, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "kafka writer initialized",
		slog.String("topic", topic),
		slog.Any("broker", cfg.Brokers))

	return w
}
