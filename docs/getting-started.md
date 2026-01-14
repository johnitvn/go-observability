# Getting Started

This page shows how to install and use the `go-observability` library in a Go service (installation,
basic initialization, and runtime usage). It is intended for service owners who want to add logging,
tracing, and metrics — not for editing the library itself.

## Install

Add the module to your project:

```bash
go get github.com/ecoma-io/go-observability
```

Or add it to your `go.mod` and run `go mod tidy`.

## Basic usage

1. Embed `BaseConfig` into your service config and load configuration (env/ldflags):

```go
type Config struct {
	observability.BaseConfig
	// your service-specific fields
}

var cfg Config
observability.LoadCfg(&cfg)
```

2. Create a logger and initialize OpenTelemetry:

```go
logger := observability.NewLogger(&cfg.BaseConfig)
defer logger.Sync()

shutdown, err := observability.InitOtel(cfg.BaseConfig)
if err != nil {
	logger.Fatal("failed to init otel", "error", err)
}
defer shutdown(context.Background())
```

3. Attach middleware or gRPC interceptors:

```go
// Gin HTTP
router := gin.New()
for _, mw := range observability.GinMiddlewareWithConfig(logger, cfg.ServiceName, nil) {
	router.Use(mw)
}

// gRPC server
srv := grpc.NewServer(
	grpc.UnaryInterceptor(observability.GrpcUnaryServerInterceptor(logger)),
)
```

## Environment variables & build metadata

- The library reads configuration from environment variables and supports build-time metadata via
  `-ldflags`.
- Common env vars: `LOG_LEVEL`, `OTEL_ENDPOINT`, `OTEL_TRACING_SAMPLE_RATE` (see `config.go` for
  full list).
- Inject service metadata at build time:

```bash
go build -ldflags "-X github.com/ecoma-io/go-observability.ServiceName=my-service -X github.com/ecoma-io/go-observability.Version=v1.2.3"
```

## Running examples

Examples are included under `examples/`. To run the Gin example locally:

```bash
cd examples/gin-service
SERVICE_NAME=my-gin-service LOG_LEVEL=debug go run main.go
# Metrics (pull mode) typically available at http://localhost:9090/metrics
```

## Metrics & Tracing at runtime

- Metrics: When configured in pull mode, a Prometheus scrape endpoint is exposed (default
  `:9090/metrics` for examples). When push mode is used, metrics are sent to the OTLP endpoint.
- Traces: Sent to the configured OTLP collector (`OTEL_ENDPOINT`) or visible through your tracing
  backend.

## Troubleshooting

- Missing traces or metrics usually mean the OTEL collector or endpoint is not reachable. Verify
  `OTEL_ENDPOINT` and network connectivity.
- If logs are not appearing, ensure `LOG_LEVEL` is set appropriately and your service calls
  `logger.Sync()` on shutdown.

## Where to go next

- API reference and integration details: see the `docs/` directory for `gin-middleware.md.md`,
  `grpc.md`, and `otel.md`.
- For development and contributions, see the repository README and CONTRIBUTING.md.
