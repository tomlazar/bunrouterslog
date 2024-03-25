package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/uptrace/bunrouter"
	"go.tomlazar.net/bunrouterslog"
)

func main() {
	router := bunrouter.New(
		bunrouter.WithMiddleware(
			bunrouterslog.NewBunrouterMiddleware(),
		),
	)

	// ... add routes to router ...

	router.GET("/.system/healthcheck", func(w http.ResponseWriter, req bunrouter.Request) error {
		logger := bunrouterslog.LoggerFromContext(req.Context())
		logger.Info("healthcheck")

		return nil
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	ctx = bunrouterslog.ContextWithLogger(ctx, logger)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("http server error", "error", err)
	}
}
