# Configuration

This project exposes a `BaseConfig` pattern for embedding into service configs. Key points:

- Priority order: LDFlags > .env file > Environment variables.
- Load via `observability.LoadCfg(&cfg)` which validates and injects build metadata.
- Modes for metrics: `IsPull()`, `IsPush()`, `IsHybrid()`.

## Example (pseudo)

```go
type Config struct {
  observability.BaseConfig
  DatabaseURL string `env:"DATABASE_URL"`
}

var cfg Config
observability.LoadCfg(&cfg)
```

## `BaseConfig` fields

The library exposes a `BaseConfig` struct with defaults and validation. Important fields:

| Field                   |                    Env var | Default          | Notes                                                         |
| ----------------------- | -------------------------: | ---------------- | ------------------------------------------------------------- |
| `ServiceName`           |             `SERVICE_NAME` | (required)       | Injected via LDFlags or env; required by `LoadCfg` validation |
| `Version`               |                          - | `dev`            | Usually injected at build-time with `-ldflags`                |
| `BuildTime`             |                          - | `unknown`        | Injected at build-time                                        |
| `LogLevel`              |                `LOG_LEVEL` | `info`           | Allowed: `debug`, `info`, `warn`, `error`                     |
| `OtelEndpoint`          |            `OTEL_ENDPOINT` | `localhost:4318` | OTLP/HTTP endpoint for traces                                 |
| `MetricsPort`           |             `METRICS_PORT` | `9090`           | HTTP port for Prometheus pull server                          |
| `OtelTracingSampleRate` | `OTEL_TRACING_SAMPLE_RATE` | `1.0`            | Trace sampling ratio (0.0 - 1.0)                              |
| `MetricsMode`           |             `METRICS_MODE` | `pull`           | `pull`, `push`, or `hybrid`                                   |
| `MetricsPath`           |             `METRICS_PATH` | `/metrics`       | Path served by Prometheus handler                             |
| `MetricsPushEndpoint`   |    `METRICS_PUSH_ENDPOINT` | -                | Required when `METRICS_MODE` is `push`/`hybrid`               |
| `MetricsPushInterval`   |    `METRICS_PUSH_INTERVAL` | `30`             | Seconds between push exports                                  |
| `MetricsProtocol`       |         `METRICS_PROTOCOL` | `http`           | `http` or `grpc` for OTLP metrics push                        |

## Validation rules performed by `LoadCfg()`

- Ensures `SERVICE_NAME` is set (or injected via LDFlags) and non-empty.
- Validates `LOG_LEVEL` is one of `debug|info|warn|error`.
- Validates `METRICS_MODE` is `pull|push|hybrid` and requires `METRICS_PUSH_ENDPOINT` for
  `push`/`hybrid`.
- Validates `METRICS_PROTOCOL` is `http` or `grpc`.

`LoadCfg` behavior summary:

- Reads `.env` file if present, otherwise reads environment variables.
- Injects build metadata via the `MetadataSetter` interface when implemented by the service config.
- Calls `finalizeAndValidate()` to enforce configuration invariants.

Examples and troubleshooting are provided in `getting-started.md` and `e2e.md`.
