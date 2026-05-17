# prom-snapshot-cli

> Command-line utility to query, export, and replay Prometheus TSDB snapshots for local debugging.

---

## Installation

```bash
go install github.com/yourusername/prom-snapshot-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/prom-snapshot-cli.git
cd prom-snapshot-cli
make build
```

---

## Usage

```bash
# Query metrics from a local TSDB snapshot
prom-snapshot-cli query --snapshot ./data/snapshots/01HXZ --metric http_requests_total

# Export snapshot data to JSON
prom-snapshot-cli export --snapshot ./data/snapshots/01HXZ --output metrics.json

# Replay a snapshot against a local Prometheus instance
prom-snapshot-cli replay --snapshot ./data/snapshots/01HXZ --target http://localhost:9090
```

### Common Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--snapshot` | Path to the TSDB snapshot directory | required |
| `--start` | Start time for query range (RFC3339) | snapshot min time |
| `--end` | End time for query range (RFC3339) | snapshot max time |
| `--output` | Output file path for exports | stdout |
| `--format` | Output format: `json`, `csv`, `text` | `text` |

---

## Requirements

- Go 1.21+
- A valid Prometheus TSDB snapshot (generated via `/api/v1/admin/tsdb/snapshot`)

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

---

## License

This project is licensed under the [MIT License](LICENSE).