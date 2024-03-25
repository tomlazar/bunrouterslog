package bunrouterslog

import (
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	inner slog.Handler
}

func NewOtelEventHandler(ctx context.Context, inner slog.Handler) slog.Handler {
	return &handler{inner: inner}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span != nil && span.IsRecording() {
		var err error

		attrs := []attribute.KeyValue{}

		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, attribute.Key(a.Key).String(a.String()))

			if a.Key != "error" {
				return false
			}
			nerr, ok := a.Value.Any().(error)
			if !ok {
				return false
			}

			err = nerr
			return false
		})

		if r.Level >= slog.LevelError {
			if err != nil {
				span.RecordError(err, trace.WithAttributes(attrs...))
			} else {
				span.RecordError(errors.New(r.Message), trace.WithAttributes(attrs...))
			}
		} else {
			span.AddEvent(r.Message, trace.WithAttributes(attrs...))
		}
	}

	return h.inner.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{inner: h.inner.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{inner: h.inner.WithGroup(name)}
}
