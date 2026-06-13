package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/industrial-iot/gateway-service/internal/forwarder"
	"github.com/industrial-iot/gateway-service/internal/mqtt"
	"github.com/industrial-iot/gateway-service/internal/parser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	logger.Info("============================================")
	logger.Info("Industrial IoT - Gateway Service")
	logger.Info("============================================")

	// --- MQTT Config (HiveMQ Cloud) ---
	mqttCfg := mqtt.Config{
		Broker:        "b99f17b791e64ce9bc41ed7cf260115e.s1.eu.hivemq.cloud",
		Port:          8883,
		ClientID:      "industrial-iot-gateway-v1",
		Username:      "hiepdh",
		Password:      "Rp_2Ngs7LQZhdb!",
		KeepAlive:     60,
		CleanSession:  false,
		AutoReconnect: true,
		Topics:        []string{"/factory/#"},
	}

	// --- Kafka Config ---
	kafkaCfg := forwarder.Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "telemetry.raw",
	}

	// --- Create components ---
	topicParser := parser.NewTopicParser(logger)
	kafkaForwarder := forwarder.NewForwarder(kafkaCfg, logger)
	defer kafkaForwarder.Close()

	// Message counter
	msgCount := 0
	startTime := time.Now()

	// --- MQTT Message Handler ---
	handler := func(topic string, payload []byte, qos byte, retain bool) {
		receivedAt := time.Now()
		msgCount++

		// Parse the message
		msg, err := topicParser.ParseMessage(topic, payload, receivedAt)
		if err != nil {
			logger.Error("Failed to parse message",
				zap.String("topic", topic),
				zap.Error(err),
			)
			return
		}
		// Print parsed message to console (pretty)
		fmt.Println("\n═══════════════════════════════════════════════════")
		fmt.Printf("📨 Message #%d\n", msgCount)
		fmt.Printf("   Topic    : %s\n", msg.Topic)
		fmt.Printf("   Factory  : %s | Area: %s | Machine: %s\n",
			msg.FactoryID, msg.Area, msg.MachineID)
		fmt.Printf("   Metric   : %s\n", msg.Metric)
		fmt.Printf("   Value    : %.4f %s\n", msg.Value, msg.Unit)
		fmt.Printf("   Quality  : %d\n", msg.Quality)
		fmt.Printf("   Time     : %s\n", msg.Timestamp.Format(time.RFC3339Nano))
		if len(msg.Tags) > 0 {
			fmt.Printf("   Tags     : %v\n", msg.Tags)
		}
		if len(msg.RawPayload) > 0 {
			fmt.Printf("   Raw JSON : %s\n", string(msg.RawPayload))
		}
		fmt.Println("═══════════════════════════════════════════════════")

		// Log structured
		logger.Info("Telemetry received",
			zap.Int("count", msgCount),
			zap.String("factory_id", msg.FactoryID),
			zap.String("area", msg.Area),
			zap.String("machine_id", msg.MachineID),
			zap.String("metric", msg.Metric),
			zap.Float64("value", msg.Value),
			zap.String("unit", msg.Unit),
		)

		// Forward to Kafka (non-blocking, with timeout)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := kafkaForwarder.Forward(ctx, msg); err != nil {
			logger.Error("Failed to forward to Kafka",
				zap.Error(err),
				zap.String("machine_id", msg.MachineID),
			)
		}

		// Print stats every 10 messages
		if msgCount%10 == 0 {
			elapsed := time.Since(startTime)
			rate := float64(msgCount) / elapsed.Seconds()
			logger.Info("Stats",
				zap.Int("total_messages", msgCount),
				zap.Float64("msg_per_second", rate),
				zap.Duration("elapsed", elapsed),
			)
		}
	}

	// --- Start MQTT Client ---
	client, err := mqtt.NewClient(mqttCfg, handler, logger)
	if err != nil {
		logger.Fatal("Failed to create MQTT client", zap.Error(err))
	}
	defer client.Disconnect(context.Background())

	logger.Info("Gateway Service is running...")
	logger.Info("Listening for MQTT messages on /factory/#")
	logger.Info("Press Ctrl+C to stop")

	// --- Wait for shutdown signal ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...",
		zap.Int("total_messages_processed", msgCount),
		zap.Duration("uptime", time.Since(startTime)),
	)
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	logger, _ := cfg.Build()
	return logger
}

// Helper: pretty print JSON
func prettyJSON(data []byte) string {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data)
	}
	pretty, err := json.MarshalIndent(obj, "   ", "  ")
	if err != nil {
		return string(data)
	}
	return string(pretty)
}
