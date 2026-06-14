package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/industrial-iot/device-service/internal/api"
	"github.com/industrial-iot/device-service/internal/repository"
	"github.com/industrial-iot/device-service/internal/service"
	"github.com/industrial-iot/device-service/internal/storage"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	logger.Info("============================================")
	logger.Info("Industrial IoT - Device Service")
	logger.Info("============================================")

	// --- PostgreSQL Config ---
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5433")
	pgDB := getEnv("POSTGRES_DATABASE", "industrial_iot")
	pgUser := getEnv("POSTGRES_USER", "iiot")
	pgPass := getEnv("POSTGRES_PASSWORD", "iiot_dev_2024")
	pgSSL := getEnv("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=25",
		pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL,
	)

	// --- Redis Config ---
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPass := getEnv("REDIS_PASSWORD", "iiot_dev_2024")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPass,
		DB:       0,
	})

	// --- MinIO Config ---
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "firmwares")
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	// --- API Port ---
	apiPort := 8083
	fmt.Sscanf(getEnv("API_PORT", "8083"), "%d", &apiPort)

	// --- Connect to PostgreSQL ---
	repo, err := repository.NewPostgresRepo(connStr, logger)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer repo.Close()

	// --- Run Migrations ---
	migrationsDir := getEnv("MIGRATIONS_DIR", "/migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "migrations" // fallback for local dev
	}

	logger.Info("Running database migrations...")
	if err := repository.RunMigrations(context.Background(), repo.Pool(), migrationsDir, logger); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// --- Connect to Redis ---
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
	} else {
		logger.Info("Connected to Redis", zap.String("addr", fmt.Sprintf("%s:%s", redisHost, redisPort)))
	}

	// --- Connect to MinIO ---
	minioCfg := storage.MinioConfig{
		Endpoint:  minioEndpoint,
		AccessKey: minioAccessKey,
		SecretKey: minioSecretKey,
		Bucket:    minioBucket,
		UseSSL:    minioUseSSL,
	}

	minioClient, err := storage.NewMinioClient(minioCfg, logger)
	if err != nil {
		logger.Warn("MinIO not available, firmware features will not work", zap.Error(err))
	}

	// --- Initialize Service Layer ---
	registryService := service.NewRegistryService(repo, logger)
	shadowService := service.NewShadowService(repo, redisClient, logger)

	// Firmware repo + service
	fwRepo := repository.NewFirmwareRepo(repo)
	var fwService *service.FirmwareService
	if minioClient != nil {
		fwService = service.NewFirmwareService(fwRepo, minioClient, logger)
	}

	// --- Start HTTP API ---
	apiServer := api.NewServer(registryService, shadowService, fwService, apiPort, logger)

	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	fmt.Println()
	logger.Info("Device Service is running...")
	logger.Info("  => HTTP API: http://localhost:" + fmt.Sprintf("%d", apiPort) + "/api/v1")
	logger.Info("  => Health:   http://localhost:" + fmt.Sprintf("%d", apiPort) + "/api/v1/health")
	logger.Info("  => Stats:    http://localhost:" + fmt.Sprintf("%d", apiPort) + "/api/v1/stats")
	logger.Info("Press Ctrl+C to stop")

	// --- Wait for shutdown signal ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
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
