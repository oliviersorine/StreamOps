package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		requestID, _ := r.Context().Value(requestIDContextKey).(string)

		s.logger.Info(
			"http request completed",
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", recorder.statusCode),
			slog.Duration("duration", time.Since(startedAt)),
		)
	})
}

func (s *Server) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")

		if requestID == "" {
			requestID = newRequestID()
		}

		w.Header().Set("X-Request-Id", requestID)

		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID, _ := r.Context().Value(requestIDContextKey).(string)

				s.logger.Error(
					"http panic recovered",
					slog.String("request_id", requestID),
					slog.Any("panic", recovered),
				)

				writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	bytes := make([]byte, 8)

	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}
