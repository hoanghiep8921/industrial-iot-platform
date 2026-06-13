package mqtt

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.uber.org/zap"
)

// MessageHandler is called for each received MQTT message
type MessageHandler func(topic string, payload []byte, qos byte, retain bool)

// Config holds MQTT client configuration
type Config struct {
	Broker        string
	Port          int
	ClientID      string
	Username      string
	Password      string
	KeepAlive     int
	CleanSession  bool
	AutoReconnect bool
	Topics        []string
}

// Client wraps the MQTT connection
type Client struct {
	cfg     Config
	cm      *mqtt.ConnectionManager
	handler MessageHandler
	logger  *zap.Logger
	mu      sync.RWMutex
}

// NewClient creates a new MQTT client
func NewClient(cfg Config, handler MessageHandler, logger *zap.Logger) (*Client, error) {
	c := &Client{
		cfg:     cfg,
		handler: handler,
		logger:  logger,
	}

	brokerURL := fmt.Sprintf("ssl://%s:%d", cfg.Broker, cfg.Port)

	// Parse broker URL
	u, err := url.Parse(brokerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid broker URL %s: %w", brokerURL, err)
	}

	logger.Info("Connecting to MQTT broker",
		zap.String("broker", brokerURL),
		zap.String("client_id", cfg.ClientID),
	)

	cliCfg := mqtt.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     uint16(cfg.KeepAlive),
		ConnectRetryDelay:             5 * time.Second,
		ConnectTimeout:                30 * time.Second,
		CleanStartOnInitialConnection: cfg.CleanSession,
		ConnectUsername:               cfg.Username,
		ConnectPassword:               []byte(cfg.Password),
		OnConnectionUp: func(cm *mqtt.ConnectionManager, connAck *paho.Connack) {
			logger.Info("MQTT connected successfully",
				zap.String("broker", brokerURL),
				zap.Uint8("reason_code", connAck.ReasonCode),
			)
			// Subscribe to topics on connection
			for _, topic := range cfg.Topics {
				suback, err := cm.Subscribe(context.Background(), &paho.Subscribe{
					Subscriptions: []paho.SubscribeOptions{
						{Topic: topic, QoS: 1},
					},
				})
				if err != nil {
					logger.Error("Failed to subscribe to topic",
						zap.String("topic", topic),
						zap.Error(err),
					)
				} else {
					logger.Info("Subscribed to topic",
						zap.String("topic", topic),
						zap.Int("reason_count", len(suback.Reasons)),
					)
				}
			}
		},
		OnConnectError: func(err error) {
			logger.Error("MQTT connection error", zap.Error(err))
		},
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.ClientID,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					go c.onMessage(pr.Packet.Topic, pr.Packet.Payload, pr.Packet.QoS, pr.Packet.Retain)
					return true, nil
				},
			},
		},
	}

	// TLS config for HiveMQ Cloud
	cliCfg.TlsCfg = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	cm, err := mqtt.NewConnection(context.Background(), cliCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MQTT connection: %w", err)
	}

	c.cm = cm
	return c, nil
}

// onMessage handles incoming MQTT messages
func (c *Client) onMessage(topic string, payload []byte, qos byte, retain bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.logger.Debug("Received MQTT message",
		zap.String("topic", topic),
		zap.Int("payload_size", len(payload)),
		zap.Int("qos", int(qos)),
	)

	if c.handler != nil {
		c.handler(topic, payload, qos, retain)
	}
}

// Disconnect gracefully closes the MQTT connection
func (c *Client) Disconnect(ctx context.Context) error {
	c.logger.Info("Disconnecting MQTT client...")
	return c.cm.Disconnect(ctx)
}
