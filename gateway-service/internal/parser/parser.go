package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// TelemetryMessage represents a parsed telemetry data point from a machine
type TelemetryMessage struct {
	// Metadata from MQTT topic: /factory/{factory_id}/{area}/{machine_id}/{metric}
	FactoryID string    `json:"factory_id"`
	Area      string    `json:"area"`
	MachineID string    `json:"machine_id"`
	Metric    string    `json:"metric"`
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"timestamp"`

	// Raw payload
	RawPayload json.RawMessage `json:"raw_payload,omitempty"`

	// Parsed fields (if JSON payload)
	Value       float64            `json:"value,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Quality     int                `json:"quality,omitempty"`
	Tags        map[string]string  `json:"tags,omitempty"`
}

// TopicParser extracts metadata from MQTT topic structure
type TopicParser struct {
	logger *zap.Logger
}

// NewTopicParser creates a new topic parser
func NewTopicParser(logger *zap.Logger) *TopicParser {
	return &TopicParser{logger: logger}
}

// ParseTopic extracts metadata from the MQTT topic
// Expected format: /factory/{factory_id}/{area}/{machine_id}/{metric}
// Examples:
//
//	/factory/HN/area_a/lathe_01/temperature
//	/factory/HCM/area_b/cnc_03/spindle_speed
//	/factory/HN/area_a/lathe_01/status
func (p *TopicParser) ParseTopic(topic string) (factoryID, area, machineID, metric string, err error) {
	// Remove leading/trailing slashes and split
	trimmed := strings.Trim(topic, "/")
	parts := strings.Split(trimmed, "/")

	if len(parts) < 5 {
		return "", "", "", "", fmt.Errorf("invalid topic format, expected /factory/{factory_id}/{area}/{machine_id}/{metric}, got: %s", topic)
	}

	// Expected: [factory, {factory_id}, {area}, {machine_id}, {metric}]
	if parts[0] != "factory" {
		return "", "", "", "", fmt.Errorf("topic must start with 'factory', got: %s", parts[0])
	}

	factoryID = parts[1]
	area = parts[2]
	machineID = parts[3]
	// metric can have sub-paths: parts[4:] joined
	metric = strings.Join(parts[4:], "/")

	return factoryID, area, machineID, metric, nil
}

// ParseMessage parses a raw MQTT message into a TelemetryMessage
func (p *TopicParser) ParseMessage(topic string, payload []byte, receivedAt time.Time) (*TelemetryMessage, error) {
	// Parse topic
	factoryID, area, machineID, metric, err := p.ParseTopic(topic)
	if err != nil {
		return nil, fmt.Errorf("topic parse error: %w", err)
	}

	msg := &TelemetryMessage{
		Topic:     topic,
		FactoryID: factoryID,
		Area:      area,
		MachineID: machineID,
		Metric:    metric,
		Timestamp: receivedAt,
		Quality:   0, // Default: good quality
	}

	// Try to parse payload as JSON
	if len(payload) > 0 {
		var payloadData map[string]interface{}
		if err := json.Unmarshal(payload, &payloadData); err == nil {
			// JSON payload - extract common fields
			msg.RawPayload = payload

			if v, ok := payloadData["value"]; ok {
				switch val := v.(type) {
				case float64:
					msg.Value = val
				case int:
					msg.Value = float64(val)
				case string:
					// For string values like status, store as-is in tags
					if msg.Tags == nil {
						msg.Tags = make(map[string]string)
					}
					msg.Tags["value_raw"] = val
				}
			}

			if u, ok := payloadData["unit"]; ok {
				if unitStr, ok := u.(string); ok {
					msg.Unit = unitStr
				}
			}

			if q, ok := payloadData["quality"]; ok {
				if qf, ok := q.(float64); ok {
					msg.Quality = int(qf)
				}
			}

			if ts, ok := payloadData["timestamp"]; ok {
				if tsStr, ok := ts.(string); ok {
					if parsed, err := time.Parse(time.RFC3339Nano, tsStr); err == nil {
						msg.Timestamp = parsed
					}
				} else if tsNum, ok := ts.(float64); ok {
					// Unix timestamp in seconds or milliseconds
					if tsNum > 1e12 {
						msg.Timestamp = time.UnixMilli(int64(tsNum))
					} else {
						msg.Timestamp = time.Unix(int64(tsNum), 0)
					}
				}
			}

			// Store all extra fields as tags
			for k, v := range payloadData {
				if k == "value" || k == "unit" || k == "quality" || k == "timestamp" {
					continue
				}
				if msg.Tags == nil {
					msg.Tags = make(map[string]string)
				}
				msg.Tags[k] = fmt.Sprintf("%v", v)
			}
		} else {
			// Non-JSON payload — try numeric value
			var numericVal float64
			if _, err := fmt.Sscanf(string(payload), "%f", &numericVal); err == nil {
				msg.Value = numericVal
			} else {
				// Store as raw string tag
				if msg.Tags == nil {
					msg.Tags = make(map[string]string)
				}
				msg.Tags["value_raw"] = string(payload)
			}
		}
	}

	p.logger.Debug("Parsed telemetry message",
		zap.String("factory_id", factoryID),
		zap.String("area", area),
		zap.String("machine_id", machineID),
		zap.String("metric", metric),
		zap.Float64("value", msg.Value),
	)

	return msg, nil
}
