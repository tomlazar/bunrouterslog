# bunrouter slog

## Bunrouter Middleware

This is a bunrouter middleware that adds a slog logger to the context, and logs some key details about the request and response. It also adds a RequestID to the context, and logs it with each request.

See the [middleware example](./examples/middleware)

## Slog -> Otel Event Handler

This is a bunrouter slog handler that sends logs to OpenTelemetry as events.

See the [otel example](./examples/sloghandler)
