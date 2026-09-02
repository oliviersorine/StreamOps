package events

import (
	"context"
	"log/slog"
)

type Router struct {
	logger       *slog.Logger
	eventChannel <-chan StreamEvent
	eventStore   *Store
}

func NewRouter(logger *slog.Logger, eventChannel <-chan StreamEvent, eventStore *Store) *Router {
	return &Router{
		logger:       logger,
		eventChannel: eventChannel,
		eventStore:   eventStore,
	}
}

func (r *Router) Start(ctx context.Context) error {
	r.logger.Info("Router started")

	for {
		select {
		case event := <-r.eventChannel:
			r.eventStore.Add(event)

			r.logger.Info(
				"Event routed to store",
				slog.String("event_id", event.ID),
				slog.String("event_type", event.Type),
			)

		case <-ctx.Done():
			r.logger.Info("Router stopping due to context cancellation")
			return nil
		}
	}
}
