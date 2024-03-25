package bunrouterslog

import (
	"log/slog"

	"github.com/uptrace/bunrouter"
)

type RequestDetector func(request bunrouter.Request) ([]any, error)

var _ RequestDetector = StandardInfoDector

func StandardInfoDector(request bunrouter.Request) ([]any, error) {
	atts := []any{
		slog.String("hostname", request.Host),
		slog.String("method", request.Method),
		slog.String("path", request.URL.Path),
		slog.String("query", request.URL.RawQuery),
		slog.String("user_agent", request.UserAgent()),
		slog.String("remote_addr", request.RemoteAddr),
	}

	if xff := request.Header.Get("X-Forwarded-For"); xff != "" {
		atts = append(atts, slog.String("x_forwarded_for", xff))
	}

	if xfp := request.Header.Get("X-Forwarded-Proto"); xfp != "" {
		atts = append(atts, slog.String("x_forwarded_proto", xfp))
	}

	if xfs := request.Header.Get("X-Forwarded-SSL"); xfs != "" {
		atts = append(atts, slog.String("x_forwarded_ssl", xfs))
	}

	return atts, nil
}

var _ RequestDetector = FlyInfoDetector

func FlyInfoDetector(request bunrouter.Request) ([]any, error) {
	if request.Header.Get("Fly-Region") == "" {
		return nil, nil
	}

	return []any{
		slog.String("fly_region", request.Header.Get("Fly-Region")),
		slog.String("fly_forwarded_port", request.Header.Get("Fly-Forwarded-Port")),
		slog.String("fly_client_ip", request.Header.Get("Fly-Client-IP")),
	}, nil
}
