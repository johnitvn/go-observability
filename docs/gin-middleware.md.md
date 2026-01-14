# Middleware & Interceptors

- Gin middleware: `observability.GinTracing(...)` — auto-creates spans, logs requests, recovers
  panics.
- gRPC interceptor: `observability.GrpcUnaryServerInterceptor(logger)` — captures status and trace
  context.

Example (Gin):

```go
r := gin.Default()
r.Use(observability.GinTracing("my-service"))
```

## Gin middleware details

The package exposes the following middleware helpers:

- `GinTracing(serviceName string)` / `GinTracingWithConfig(serviceName, cfg)` — starts server spans
  for incoming HTTP requests and extracts W3C trace context from headers. Injects `X-Trace-ID` in
  responses when available.
- `GinLogger(logger *Logger)` / `GinLoggerWithConfig(logger, cfg)` — logs requests with latency,
  status, client IP, user-agent and trace context.
- `GinRecovery(logger *Logger)` / `GinRecoveryWithConfig(logger, cfg)` — recovers panics, logs stack
  trace and returns structured `ErrorResponse` JSON with optional `trace_id`.
- `GinMiddleware(logger, serviceName)` — convenience to return the full chain: tracing, recovery,
  logger.

`ObservabilityMiddlewareConfig` supports:

- `ExcludedPaths []string` — exact-match paths to skip (e.g., `/health`, `/metrics`).
- `SkipRoute func(path string) bool` — user-supplied predicate to decide skipping (takes precedence
  over `ExcludedPaths`).

Usage notes:

- Skipped routes do not create spans, do not log, and do not record metrics.
- `GinRecovery` logs the full stack trace and returns a consistent JSON error structure for clients.

## gRPC interceptors

- `GrpcUnaryServerInterceptor(logger)` — logs unary RPCs with gRPC status, latency and trace
  context.
- `GrpcStreamServerInterceptor(logger)` — logs streaming RPCs with similar fields.
- `GrpcUnaryRecoveryInterceptor` and `GrpcStreamRecoveryInterceptor` — recover from panics in
  handlers and return `codes.Internal`.
- `GrpcUnaryInterceptors(logger)` and `GrpcStreamInterceptors(logger)` — helper to return
  interceptor chains (recovery + logging).

Usage example:

```go
server := grpc.NewServer(
  grpc.ChainUnaryInterceptor(observability.GrpcUnaryInterceptors(logger)...),
  grpc.ChainStreamInterceptor(observability.GrpcStreamInterceptors(logger)...),
)
```

## Notes from code review

- gRPC recovery interceptors inject `trace_id` into trailers when available to aid debugging.
- `setTrailer` wraps `grpc.SetTrailer` and deliberately ignores the returned values to satisfy
  linters; this is acceptable but documented for reviewers.
