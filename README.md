# Polkadot Validator Rewards API

A Go REST API that fetches and aggregates staking rewards for Polkadot validators using Subscan API.

## Quick Start

```bash
go mod download
make run
```

Server runs on `http://localhost:8080`

## Environment Variables

Copy the example file and configure:

```bash
cp .env.example .env
```

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `ENVIRONMENT` | development | Environment name |
| `LOG_LEVEL` | info | Log level (debug, info, warn, error) |
| `SUBSCAN_BASE_URL` | https://polkadot.api.subscan.io | Subscan API URL |
| `SUBSCAN_API_KEY` | - | Optional API key (rate limit: 5 req/s without) |
| `MAX_PAGES` | 10 | Max pages to fetch |
| `ROWS_PER_PAGE` | 100 | Results per page |
| `HTTP_TIMEOUT_SECONDS` | 30 | HTTP client timeout |
| `BLOCKS_PER_ERA` | 14400 | Polkadot blocks per era |
| `PLANCK_TO_DOT` | 10000000000 | Planck to DOT conversion |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/validators/:address/rewards/eras` | Rewards aggregated by era |
| GET | `/validators/:address/rewards/daily` | Rewards aggregated by day |
| GET | `/validators/:address/rewards/monthly` | Rewards aggregated by month |

## Example

```bash
curl http://localhost:8080/validators/14ShUZUYUR35RBZW6uVVt1zXDxmSQddkeDdXf1JkMA6P721N/rewards/eras
```

Response:
```json
{
    "validator": "14ShUZUYUR35RBZW6uVVt1zXDxmSQddkeDdXf1JkMA6P721N",
    "rewards": [
        { "era": 450, "amount": "173.4499" },
        { "era": 297, "amount": "90.6556" }
    ]
}
```

## Commands

```bash
make build             # Build binary
make run               # Run server
make test              # Run unit tests
make test-cover        # Run unit tests with coverage
make test-integration  # Run all tests (including integration)
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make clean             # Clean build artifacts
```

## Project Structure

```
.
├── main.go              # Entry point with graceful shutdown
├── config/              # Configuration loading
├── controllers/         # HTTP handlers
├── services/            # External API clients (Subscan)
├── helpers/             # Business logic (aggregation)
├── routes/              # Route definitions
├── types/               # Data structures
└── tests/               # Integration tests
```

## Testing

Unit tests use mocked dependencies and table-driven tests:

```bash
make test-cover
```

Coverage:
- `controllers`: 80%+
- `helpers`: 100%

## Docker

```bash
make docker-build
make docker-run
```

The container runs as non-root user and includes health checks.
