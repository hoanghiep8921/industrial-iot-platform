package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/industrial-iot/telemetry-service/internal/api"
	"github.com/industrial-iot/telemetry-service/internal/consumer"
	"github.com/industrial-iot/telemetry-service/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	logger.Info("============================================")
	logger.Info("Industrial IoT - Telemetry Service")
	logger.Info("============================================")

	// --- TimescaleDB Config ---
	dbCfg := storage.Config{
		Host:           getEnv("TSDB_HOST", "localhost"),
		Port:           5432,
		Database:       getEnv("TSDB_DATABASE", "telemetry"),
		Username:       getEnv("TSDB_USER", "iiot"),
		Password:       getEnv("TSDB_PASSWORD", "iiot_dev_2024"),
		SSLMode:        getEnv("TSDB_SSLMODE", "disable"),
		MaxConnections: 25,
	}

	// --- Kafka Config ---
	kafkaCfg := consumer.Config{
		Brokers:       []string{getEnv("KAFKA_BROKER", "localhost:9092")},
		Topic:         getEnv("KAFKA_TOPIC", "telemetry.raw"),
		ConsumerGroup: "telemetry-service",
		NumConsumers:  4,
	}

	// --- Connect to TimescaleDB ---
	db, err := storage.NewTimescaleDB(dbCfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to TimescaleDB", zap.Error(err))
	}
	defer db.Close()

	// --- Create Kafka consumer with handler ---
	msgCount := 0
	startTime := time.Now()

	handler := func(msg *consumer.TelemetryMessage) error {
		msgCount++

		record := &storage.TelemetryRecord{
			Time:       msg.Timestamp,
			DeviceID:   msg.MachineID,
			FactoryID:  msg.FactoryID,
			Area:       msg.Area,
			MetricName: msg.Metric,
			Value:      msg.Value,
			Unit:       msg.Unit,
			Quality:    msg.Quality,
			Tags:       msg.Tags,
			RawPayload: msg.RawPayload,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.Insert(ctx, record); err != nil {
			return err
		}

		// Log progress every 100 messages
		if msgCount%100 == 0 {
			elapsed := time.Since(startTime)
			rate := float64(msgCount) / elapsed.Seconds()
			logger.Info("Processing stats",
				zap.Int("stored_messages", msgCount),
				zap.Float64("msg_per_second", rate),
			)
		}

		return nil
	}

	kafkaConsumer := consumer.NewConsumer(kafkaCfg, handler, logger)

	// --- Start Kafka consumer in background ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println()
	logger.Info("Starting Kafka consumer...")
	if err := kafkaConsumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start Kafka consumer", zap.Error(err))
	}

	// --- Start HTTP API ---
	apiServer := api.NewServer(db, 8082, logger)

	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	logger.Info("Telemetry Service is running...")
	logger.Info("  => HTTP API: http://localhost:8082/api/v1")
	logger.Info("  => Health:   http://localhost:8082/api/v1/health")
	logger.Info("  => Recent:   http://localhost:8082/api/v1/telemetry/recent?limit=50")
	logger.Info("Press Ctrl+C to stop")

	// --- Wait for shutdown signal ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...",
		zap.Int("total_messages_stored", msgCount),
		zap.Duration("uptime", time.Since(startTime)),
	)

	cancel()
	kafkaConsumer.Close()
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	logger, _ := cfg.Build()
	return logger
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
