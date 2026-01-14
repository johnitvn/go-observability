# Logging

The library provides a Zap-based logger with service metadata auto-attached.

Create a logger:

```go
logger := observability.NewLogger(&cfg.BaseConfig)
defer logger.Sync()
logger.Info("starting service", "version", observability.GetVersion())
```

Logs are JSON to stdout and include service metadata for easy ingestion.

## Implementation details

`NewLogger` constructs a Zap `SugaredLogger` with JSON encoder and the following characteristics:

- Time encoding: `ISO8601` (field key `timestamp`).
- Output: `os.Stdout` (JSON lines suitable for log collectors).
- Caller information and stacktraces included for error level logs.
- Pre-attaches `service` and `version` fields from `BaseConfig`.

The `Logger` wrapper exposes convenience methods: `Info`, `Error`, `Debug`, `Warn`, `Fatal`, `Sync`.

## Best practices

- Always call `defer logger.Sync()` to flush any buffered logs before process exit.
- Use structured key/value pairs, e.g., `logger.Info("cache miss", "key", key)`.
- Prefer the wrapper methods on `observability.Logger` to keep log format consistent across
  services.

## Notes from code review

- `Sync()` error is intentionally ignored in the implementation to avoid noisy shutdown errors;
  consider logging the error when debugging shutdown issues.
