package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/oliviersorine/StreamOps/internal/config"
	"github.com/oliviersorine/StreamOps/internal/events"
)

type Server struct {
	cfg        config.HTTPConfig
	logger     *slog.Logger
	eventStore *events.Store
	eventBus   *events.Bus
	version    string
}

func NewServer(
	cfg config.HTTPConfig,
	logger *slog.Logger,
	eventStore *events.Store,
	eventBus *events.Bus,
	version string,
) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		eventStore: eventStore,
		eventBus:   eventBus,
		version:    version,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /version", s.handleVersion)
	mux.HandleFunc("GET /time", s.handleTime)
	mux.HandleFunc("GET /events", s.handleListEvents)
	mux.HandleFunc("POST /events", s.handlePostEvent)

	server := &http.Server{
		Addr:              s.address(),
		Handler:           s.recoveryMiddleware(s.requestIDMiddleware(s.loggingMiddleware(mux))),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("Server is up and running", slog.String("address", s.address()))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down server is requested")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		s.logger.Info("Server has been shut down")

		return nil

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("Start server error: %w", err)
		}

		return nil
	}
}

func (s *Server) address() string {
	return fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": s.version})
}

func (s *Server) handleTime(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"time": time.Now().UTC().Format(time.RFC3339)})
}

func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"events": s.eventStore.List(),
	})
}

type eventRequest struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (s *Server) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	var request eventRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Failed to parse request body")
		return
	}

	if request.Type == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "Event type is required")
		return
	}

	if request.Payload == nil {
		request.Payload = map[string]any{}
	}

	event := events.New("test", request.Type, request.Payload)

	err := s.eventBus.Publish(r.Context(), event)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "publish_error", fmt.Sprintf("Failed to publish event: %v", err))
		return
	}

	s.logger.Info(
		"Event published",
		slog.String("event_id", event.ID),
		slog.String("event_type", event.Type),
	)

	writeJSON(w, http.StatusCreated, map[string]any{"event": event})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("failed to write JSON response: %v", err), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, map[string]string{"error": code, "message": message})
}
