package obs

import "context"

type Scene struct {
	Name string `json:"name"`
}

type Client interface {
	Connect(ctx context.Context) error
	IsConnected() bool
	ListScenes(ctx context.Context) ([]Scene, error)
	SwitchScene(ctx context.Context, sceneName string) error
	Close(ctx context.Context) error
}
