package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/industrial-iot/alarm-service/internal/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Config struct {
	Brokers []string
	Topic   string
}

type Publisher struct {
	cfg    Config
	writer *kafka.Writer
	logger *zap.Logger
}

func NewPublisher(cfg Config, logger *zap.Logger) *Publisher {
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

	return &Publisher{
		cfg:    cfg,
		writer: writer,
		logger: logger,
	}
}

func (f *Publisher) PublishAlarmEvent(ctx context.Context, eventType string, alarm *model.Alarm) error {
	event := model.AlarmEvent{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Alarm:     alarm,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal alarm event: %w", err)
	}

	err = f.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(alarm.MachineID), // partition by machine_id to preserve order
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		f.logger.Error("Failed to publish alarm event to Kafka",
			zap.Error(err),
			zap.String("topic", f.cfg.Topic),
			zap.String("event_type", eventType),
		)
		return fmt.Errorf("kafka write error: %w", err)
	}

	f.logger.Info("Alarm event published to Kafka",
		zap.String("topic", f.cfg.Topic),
		zap.String("event_type", eventType),
		zap.String("alarm_id", alarm.ID),
	)

	return nil
}

func (f *Publisher) Close() error {
	f.logger.Info("Closing Kafka publisher...")
	return f.writer.Close()
}
