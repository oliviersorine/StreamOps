package events

import (
	"context"
)

type Bus struct {
	ch chan StreamEvent
}

func NewBus(bufferSize int) *Bus {
	return &Bus{
		ch: make(chan StreamEvent, bufferSize),
	}
}

func (b *Bus) Publish(ctx context.Context, event StreamEvent) error {
	select {
	case b.ch <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Bus) Events() <-chan StreamEvent {
	return b.ch
}
