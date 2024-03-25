package bunrouterslog

import (
	"context"
	"log/slog"
)

type bunrouterSlogKey struct{}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	var (
		logger *slog.Logger
		ok     bool

		val = ctx.Value(bunrouterSlogKey{})
	)
	if val == nil {
		logger = slog.Default()
		logger.Warn("no logger in context")
	} else {
		logger, ok = val.(*slog.Logger)
		if !ok {
			logger = slog.Default()
			logger.Warn("logger in context is not of type *slog.Logger")
		}
	}

	return logger
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, bunrouterSlogKey{}, logger)
}
