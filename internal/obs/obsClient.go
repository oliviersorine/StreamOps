package obs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/andreykaipov/goobs"
	"github.com/oliviersorine/StreamOps/internal/config"
)

type WebSocketClient struct {
	logger *slog.Logger
	cfg    config.OBSConfig

	client *goobs.Client
	connected bool
}

func NewWebSocketClient(logger *slog.Logger, cfg config.OBSConfig) *WebSocketClient {
	return &WebSocketClient{
		logger: logger,
		cfg:    cfg,
	}
}

func (c *WebSocketClient) Connect(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	client, err := goobs.New(address, goobs.WithPassword(c.cfg.Password))
	if err != nil {
		c.connected = false
		return fmt.Errorf("failed to connect to OBS WebSocket: %w", err)
	}

	c.client = client
	c.connected = true

	c.logger.Info(
		"Connected to OBS WebSocket",
		slog.String("address", address),
	)

	return nil
}

func (c *WebSocketClient) IsConnected() bool {
	return c.connected && c.client != nil
}

func (c *WebSocketClient) ListScenes(ctx context.Context) ([]Scene, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to OBS WebSocket")
	}

	response, err := c.client.Scenes.GetSceneList()
	if err != nil {
		return nil, fmt.Errorf("failed to list scenes: %w", err)
	}

	scenes := make([]Scene, len(response.Scenes))
	for i, s := range response.Scenes {
		scenes[i] = Scene{
			Name: s.SceneName,
		}
	}
	return scenes, nil
}

func (c *WebSocketClient) SwitchScene(ctx context.Context, sceneName string) error {
	return nil
}

func (c *WebSocketClient) Close(ctx context.Context) error {
	return nil
}
