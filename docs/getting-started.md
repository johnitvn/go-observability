# Getting Started

Quickstart (local):

Install dependencies and run docs preview:

```bash
python -m pip install -r requirements.txt
mkdocs serve -a 0.0.0.0:8000
```

Run an example service (from repository root):

```bash
cd examples/gin-service
go run main.go
# Metrics available at http://localhost:9090/metrics (if OTEL in pull mode)
```
