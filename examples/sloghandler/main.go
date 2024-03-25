package main

import (
	"context"
	"log/slog"
	"os"

	"go.tomlazar.net/bunrouterslog"
)

func main() {
	otelHandler := bunrouterslog.NewOtelEventHandler(
		context.Background(),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	logger := slog.New(otelHandler)
	logger.Info("application start", "version", "1.0.0")
}
