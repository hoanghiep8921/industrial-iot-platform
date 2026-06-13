package forwarder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Config holds Kafka forwarder configuration
type Config struct {
	Brokers []string
	Topic   string
}

// Forwarder sends parsed telemetry data to Kafka
type Forwarder struct {
	cfg    Config
	writer *kafka.Writer
	logger *zap.Logger
}

// NewForwarder creates a new Kafka forwarder
func NewForwarder(cfg Config, logger *zap.Logger) *Forwarder {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.Hash{},
		BatchSize:    100,
		BatchTimeout: 500 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
		Compression:  kafka.Snappy,
	}

	return &Forwarder{
		cfg:    cfg,
		writer: writer,
		logger: logger,
	}
}

// Forward sends a telemetry message to Kafka
func (f *Forwarder) Forward(ctx context.Context, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = f.writer.WriteMessages(ctx, kafka.Message{
		Key:   nil,
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		f.logger.Error("Failed to forward message to Kafka",
			zap.Error(err),
			zap.String("topic", f.cfg.Topic),
		)
		return fmt.Errorf("kafka write error: %w", err)
	}

	f.logger.Debug("Message forwarded to Kafka",
		zap.String("topic", f.cfg.Topic),
		zap.Int("size", len(data)),
	)

	return nil
}

// Close gracefully closes the Kafka writer
func (f *Forwarder) Close() error {
	f.logger.Info("Closing Kafka forwarder...")
	return f.writer.Close()
}
