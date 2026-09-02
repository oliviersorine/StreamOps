package main

import (
	"context"
	"fmt"
	"github.com/oliviersorine/StreamOps/internal/api"
	"github.com/oliviersorine/StreamOps/internal/config"
	"github.com/oliviersorine/StreamOps/internal/events"
	"github.com/oliviersorine/StreamOps/internal/logging"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	command := os.Args[1]

	switch command {
	case "start":
		if err := runStart(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println(version)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %v\n", command)
		printUsage()
		os.Exit(2)
	}
}

func runStart() error {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		return err
	}

	logger := logging.New(cfg.Logging)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	logger.Info(
		"StreamOps started",
		slog.String("app", cfg.App.Name),
		slog.String("environment", cfg.App.Environment),
		slog.String("http_host", cfg.HTTP.Host),
		slog.Int("http_port", cfg.HTTP.Port),
	)

	eventStore := events.NewStore(100)
	eventBus := events.NewBus(100)

	eventRoute := events.NewRouter(
		logger,
		eventBus.Events(),
		eventStore,
	)

	go func() {
		if err := eventRoute.Start(ctx); err != nil {
			logger.Error("Error starting event route", slog.Any("error", err))
		}
	}()

	apiServer := api.NewServer(cfg.HTTP, logger, eventStore, eventBus, version)
	if err := apiServer.Start(ctx); err != nil {
		return err
	}

	logger.Info("StreamOps stopped")

	return nil
}

func printUsage() {
	fmt.Println(`Usage:
  streamops start
  streamops version`)
}
