package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// TelemetryMessage matches the structure from gateway-service parser
type TelemetryMessage struct {
	FactoryID  string            `json:"factory_id"`
	Area       string            `json:"area"`
	MachineID  string            `json:"machine_id"`
	Metric     string            `json:"metric"`
	Topic      string            `json:"topic"`
	Timestamp  time.Time         `json:"timestamp"`
	RawPayload json.RawMessage   `json:"raw_payload,omitempty"`
	Value      float64           `json:"value,omitempty"`
	Unit       string            `json:"unit,omitempty"`
	Quality    int               `json:"quality,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
}

// MessageHandler is called for each consumed message
type MessageHandler func(msg *TelemetryMessage) error

// Config holds Kafka consumer configuration
type Config struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	NumConsumers  int
}

// Consumer reads telemetry messages from Kafka
type Consumer struct {
	cfg     Config
	handler MessageHandler
	logger  *zap.Logger
	readers []*kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg Config, handler MessageHandler, logger *zap.Logger) *Consumer {
	return &Consumer{
		cfg:     cfg,
		handler: handler,
		logger:  logger,
	}
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka consumer",
		zap.String("topic", c.cfg.Topic),
		zap.String("group", c.cfg.ConsumerGroup),
		zap.Int("consumers", c.cfg.NumConsumers),
	)

	for i := 0; i < c.cfg.NumConsumers; i++ {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        c.cfg.Brokers,
			Topic:          c.cfg.Topic,
			GroupID:        c.cfg.ConsumerGroup,
			MinBytes:       1,
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
		})

		c.readers = append(c.readers, reader)

		go func(id int, r *kafka.Reader) {
			c.logger.Info("Consumer worker started", zap.Int("worker_id", id))
			for {
				select {
				case <-ctx.Done():
					c.logger.Info("Consumer worker stopping", zap.Int("worker_id", id))
					return
				default:
					msg, err := r.FetchMessage(ctx)
					if err != nil {
						if ctx.Err() != nil {
							return
						}
						c.logger.Error("Failed to fetch message",
							zap.Int("worker_id", id),
							zap.Error(err),
						)
						time.Sleep(time.Second)
						continue
					}

					// Parse message
					var telMsg TelemetryMessage
					if err := json.Unmarshal(msg.Value, &telMsg); err != nil {
						c.logger.Error("Failed to unmarshal message",
							zap.Int("worker_id", id),
							zap.Error(err),
							zap.String("raw", string(msg.Value)),
						)
						r.CommitMessages(ctx, msg)
						continue
					}

					// Process message
					if err := c.handler(&telMsg); err != nil {
						c.logger.Error("Failed to process message",
							zap.Int("worker_id", id),
							zap.Error(err),
							zap.String("machine_id", telMsg.MachineID),
						)
					}

					// Commit offset
					if err := r.CommitMessages(ctx, msg); err != nil {
						c.logger.Error("Failed to commit offset",
							zap.Int("worker_id", id),
							zap.Error(err),
						)
					}
				}
			}
		}(i, reader)
	}

	return nil
}

// Close gracefully shuts down all consumers
func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka consumers...")
	for i, reader := range c.readers {
		if err := reader.Close(); err != nil {
			return fmt.Errorf("failed to close reader %d: %w", i, err)
		}
	}
	return nil
}
