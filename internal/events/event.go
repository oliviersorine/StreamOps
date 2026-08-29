package events

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type StreamEvent struct {
	ID	string	`json:"id"`
	Source	string	`json:"source"`
	Type	string	`json:"type"`
	Payload	map[string]any	`json:"payload"`
	CreatedAt	time.Time	`json:"created_at"`
}

func New(source string, eventType string, payload map[string]any) StreamEvent {
	return StreamEvent{
		ID: newID(),
		Source: source,
		Type: eventType,
		Payload: payload,
		CreatedAt: time.Now().UTC(),
	}
}

func newID() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}