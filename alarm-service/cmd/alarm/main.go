package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/industrial-iot/alarm-service/internal/api"
	"github.com/industrial-iot/alarm-service/internal/consumer"
	"github.com/industrial-iot/alarm-service/internal/publisher"
	"github.com/industrial-iot/alarm-service/internal/rules"
	"github.com/industrial-iot/alarm-service/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := initLogger()
	defer logger.Sync()

	logger.Info("============================================")
	logger.Info("Industrial IoT - Alarm Service")
	logger.Info("============================================")

	// --- Cấu hình PostgreSQL ---
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPortStr := getEnv("POSTGRES_PORT", "5433")
	pgPort, err := strconv.Atoi(pgPortStr)
	if err != nil {
		pgPort = 5433
	}
	pgDb := getEnv("POSTGRES_DATABASE", "industrial_iot")
	pgUser := getEnv("POSTGRES_USER", "iiot")
	pgPass := getEnv("POSTGRES_PASSWORD", "iiot_dev_2024")

	pgConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgDb)

	// --- Cấu hình Redis ---
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPortStr := getEnv("REDIS_PORT", "6379")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		redisPort = 6379
	}
	redisPass := getEnv("REDIS_PASSWORD", "iiot_dev_2024")

	// --- Cấu hình Kafka ---
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	telemetryTopic := getEnv("KAFKA_TELEMETRY_TOPIC", "telemetry.raw")
	alarmTopic := getEnv("KAFKA_ALARM_TOPIC", "alarm.events")

	// --- Cấu hình HTTP Port ---
	apiPortStr := getEnv("API_PORT", "8086")
	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		apiPort = 8086
	}

	// 1. Kết nối PostgreSQL
	pgRepo, err := storage.NewPostgresRepo(pgConnStr, logger)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer pgRepo.Close()

	// 2. Chạy Database Migrations
	logger.Info("Running database migrations...")
	migrationsDir := getEnv("MIGRATIONS_DIR", "./migrations")
	if err := storage.RunMigrations(context.Background(), pgRepo.Pool(), migrationsDir, logger); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// 3. Kết nối Redis
	redisCache, err := storage.NewRedisCache(redisHost, redisPort, redisPass, logger)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisCache.Close()

	// 4. Khởi tạo Kafka Publisher
	logger.Info("Initializing Kafka publisher...")
	pubCfg := publisher.Config{
		Brokers: []string{kafkaBroker},
		Topic:   alarmTopic,
	}
	pub := publisher.NewPublisher(pubCfg, logger)
	defer pub.Close()

	// 5. Khởi tạo Rules Evaluator
	evaluator := rules.NewEvaluator(pgRepo, redisCache, pub, logger)

	// 6. Khởi tạo & Chạy Kafka Consumer
	logger.Info("Initializing Kafka consumer...")
	consumerCfg := consumer.Config{
		Brokers:       []string{kafkaBroker},
		Topic:         telemetryTopic,
		ConsumerGroup: "alarm-service",
		NumConsumers:  4,
	}

	handler := func(ctx context.Context, msg *consumer.TelemetryMessage) error {
		return evaluator.Evaluate(ctx, msg)
	}

	kafkaConsumer := consumer.NewConsumer(consumerCfg, handler, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := kafkaConsumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start Kafka consumer", zap.Error(err))
	}

	// 7. Khởi tạo & Chạy HTTP API Server
	apiServer := api.NewServer(pgRepo, apiPort, logger)
	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Fatal("HTTP API server failed", zap.Error(err))
		}
	}()

	logger.Info("Alarm Service is running...")
	logger.Info(fmt.Sprintf("  => HTTP API: http://localhost:%d/api/v1", apiPort))
	logger.Info(fmt.Sprintf("  => Health:   http://localhost:%d/api/v1/health", apiPort))
	logger.Info("Press Ctrl+C to stop")

	// 8. Đợi tín hiệu dừng chương trình
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down Alarm Service...")
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
