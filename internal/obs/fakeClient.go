package obs

import (
	"context"
	"fmt"
)

type FakeClient struct {
	connected    bool
	scenes       []Scene
	currentScene string
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		connected: false,
		scenes: []Scene{
			{Name: "Scene 1"},
			{Name: "Scene 2"},
		},
	}
}

func (c *FakeClient) Connect(ctx context.Context) error {
	fmt.Println("Fake OBS client connected")
	c.connected = true
	return nil
}

func (c *FakeClient) IsConnected() bool {
	return c.connected
}

func (c *FakeClient) ListScenes(ctx context.Context) ([]Scene, error) {
	return []Scene{
		{Name: "Scene 1"},
		{Name: "Scene 2"},
	}, nil
}

func (c *FakeClient) SwitchScene(ctx context.Context, sceneName string) error {
	if !c.connected {
		return fmt.Errorf("OBS client is not connected")
	}

	for _, scene := range c.scenes {
		if scene.Name == sceneName {
			c.currentScene = sceneName
			return nil
		}
	}

	return fmt.Errorf("scene not found: %s", sceneName)
}

func (c *FakeClient) Close(ctx context.Context) error {
	c.connected = false
	return nil
}
