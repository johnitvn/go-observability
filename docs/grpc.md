# gRPC Interceptors

This project provides gRPC interceptors for unary and streaming RPCs that add logging, tracing, and
panic recovery.

Key helpers:

- `GrpcUnaryServerInterceptor(logger *observability.Logger)` — logs unary RPC calls with latency,
  gRPC status, and trace context.
- `GrpcStreamServerInterceptor(logger *observability.Logger)` — logs streaming RPCs and their
  metadata.
- `GrpcUnaryRecoveryInterceptor(logger *observability.Logger)` — recovers from panics and converts
  them to `codes.Internal` errors, attaching `trace_id` metadata when available.
- `GrpcStreamRecoveryInterceptor(logger *observability.Logger)` — similar to unary recovery but for
  streams.
- `GrpcUnaryInterceptors(logger)` / `GrpcStreamInterceptors(logger)` — return interceptor chains
  (recovery + logging) for easy wiring.

Usage example:

```go
server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(observability.GrpcUnaryInterceptors(logger)...),
    grpc.ChainStreamInterceptor(observability.GrpcStreamInterceptors(logger)...),
)
```

Notes:

- Recovery interceptors log panics with stack traces and attempt to set `trace_id` in response
  trailers for easier debugging.
- Interceptors rely on OpenTelemetry context propagation; ensure `InitOtel` has been called in your
  service initialization.
