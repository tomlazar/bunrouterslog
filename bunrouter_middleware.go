package bunrouterslog

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"
)

type ResponseInfo struct {
	Status       string
	BytesWritten int
}

type bunrouterMiddleware struct {
	detectors []RequestDetector
}

func WithRequestDetectors(detectors ...RequestDetector) Option {
	return func(m *bunrouterMiddleware) {
		m.detectors = append(m.detectors, detectors...)
	}
}

type Option func(m *bunrouterMiddleware)

func NewBunrouterMiddleware(options ...Option) bunrouter.MiddlewareFunc {
	m := &bunrouterMiddleware{
		detectors: []RequestDetector{
			StandardInfoDector,
		},
	}
	for _, opt := range options {
		opt(m)
	}

	return m.Middleware
}

func (m *bunrouterMiddleware) Middleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		var (
			logger  = LoggerFromContext(r.Context())
			id, err = NewRequestID()
		)
		if err != nil {
			return fmt.Errorf("failed to generate request id: %w", err)
		}

		logger = logger.With(slog.String("http_request_id", id.String()))
		var (
			ctx   = ContextWithRequestID(ContextWithLogger(r.Context(), logger), id)
			attrs = []any{
				slog.String("request_id", id.String()),
			}
		)

		for _, detector := range m.detectors {
			nattr, err := detector(r)
			if err != nil {
				return fmt.Errorf("failed to detect request info: %w", err)
			}
			attrs = append(attrs, nattr...)
		}

		logger.Info("request started", attrs...)

		var (
			now = time.Now()
			rw  = wrap(w)
		)
		err = next(rw, r.WithContext(ctx))
		if err != nil {
			logger.Error("request failed",
				slog.String("error", err.Error()),
				slog.Int("http_status_code", rw.statusCode),
				slog.Int("http_bytes_written", rw.bytesWritten),
				slog.Duration("http_duration", time.Since(now)),
			)
		} else {
			logger.Info("request finished",
				slog.Int("http_status_code", rw.statusCode),
				slog.Int("http_bytes_written", rw.bytesWritten),
				slog.Duration("http_duration", time.Since(now)),
			)
		}

		return err
	}
}
