package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type TelemetryMessage struct {
	FactoryID string            `json:"factory_id"`
	Area      string            `json:"area"`
	MachineID string            `json:"machine_id"`
	Metric    string            `json:"metric"`
	Topic     string            `json:"topic"`
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit,omitempty"`
	Quality   int               `json:"quality"`
}

func main() {
	log.Println("==========================================================")
	log.Println("Industrial IoT - Laser Machine Telemetry Simulator")
	log.Println("==========================================================")

	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TELEMETRY_TOPIC", "telemetry.raw")
	machineID := "LASER-3015-3000W-001"
	factoryID := "factory-01"
	area := "cutting"

	log.Printf("Connecting to Kafka at %s (Topic: %s)...", kafkaBroker, topic)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	ctx := context.Background()
	cycle := 0
	mode := "NORMAL" // NORMAL -> TEMPSPICK -> NORMAL -> GASDROP -> NORMAL

	log.Println("Simulator started! Press Ctrl+C to exit.")
	log.Println("----------------------------------------------------------")

	for {
		cycle++
		// Determine simulation mode based on cycle counts
		// Normal for 4 cycles, Temp spike for 3 cycles, Normal for 4 cycles, Gas pressure drop for 3 cycles...
		modCycle := cycle % 14
		if modCycle >= 1 && modCycle <= 4 {
			mode = "NORMAL"
		} else if modCycle >= 5 && modCycle <= 7 {
			mode = "TEMPSPICK"
		} else if modCycle >= 8 && modCycle <= 11 {
			mode = "NORMAL"
		} else {
			mode = "GASDROP"
		}

		log.Printf("[Cycle %d] Mode: %s\n", cycle, mode)

		// 1. Laser Power (Normal ~2850W, off when idle or load change)
		power := 2800.0 + rand.Float64()*100.0

		// 2. Laser Temperature (Normal 42-45°C, Spikes to 58-60°C in TEMPSPICK)
		temp := 42.0 + rand.Float64()*3.0
		if mode == "TEMPSPICK" {
			temp = 57.0 + rand.Float64()*3.0 // triggers rule > 55°C
		}

		// 3. Gas Pressure (Normal 12-14 bar, drops to 3.0-3.5 bar in GASDROP)
		pressure := 12.0 + rand.Float64()*2.0
		if mode == "GASDROP" {
			pressure = 3.0 + rand.Float64()*0.8 // triggers rule < 4.0 bar
		}

		// 4. Cutting Speed (Normal 14-16 m/min)
		speed := 14.0 + rand.Float64()*2.0

		// 5. Nozzle Height (Normal 0.75-0.85 mm)
		nozzle := 0.75 + rand.Float64()*0.1

		// Prepare messages
		metrics := []struct {
			name  string
			value float64
			unit  string
		}{
			{"laser_power", power, "W"},
			{"laser_temp", temp, "°C"},
			{"gas_pressure", pressure, "bar"},
			{"cutting_speed", speed, "m/min"},
			{"nozzle_height", nozzle, "mm"},
		}

		var messages []kafka.Message
		for _, m := range metrics {
			payload := TelemetryMessage{
				FactoryID: factoryID,
				Area:      area,
				MachineID: machineID,
				Metric:    m.name,
				Topic:     fmt.Sprintf("devices/%s/telemetry", machineID),
				Timestamp: time.Now(),
				Value:     m.value,
				Unit:      m.unit,
				Quality:   0, // good
			}

			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Failed to marshal payload for %s: %v\n", m.name, err)
				continue
			}

			messages = append(messages, kafka.Message{
				Key:   []byte(machineID),
				Value: data,
			})
			log.Printf("  -> Telemetry: %s = %.2f %s\n", m.name, m.value, m.unit)
		}

		// Write to Kafka
		err := writer.WriteMessages(ctx, messages...)
		if err != nil {
			log.Printf("Error sending telemetry to Kafka: %v\n", err)
		} else {
			log.Println("  Successfully sent metrics to Kafka.")
		}

		time.Sleep(3 * time.Second)
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
