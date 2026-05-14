# go-template

A template repository for building Web APIs in Go with Echo, GORM, and PostgreSQL, following Clean Architecture. It ships with a Dev Container / Docker Compose setup, `golang-migrate` for schema migrations, `air` for hot reload, and a layered test suite with mocks.

## Tech Stack

- Go 1.26
- [Echo v4](https://echo.labstack.com/) — HTTP framework
- [GORM](https://gorm.io/) + PostgreSQL 17
- [golang-migrate](https://github.com/golang-migrate/migrate) — schema migrations
- [air](https://github.com/air-verse/air) — hot reload
- [testify](https://github.com/stretchr/testify) — assertions and mocks
- Dev Container (VS Code) + Docker Compose

## Project Structure

```
.
├── app/
│   ├── cmd/                  # Entry point (main.go)
│   ├── infrastructure/
│   │   ├── config/           # Environment-based configuration loader
│   │   └── db/               # DB connection setup / teardown
│   ├── internal/
│   │   ├── entity/           # Domain models
│   │   ├── repository/       # Persistence layer
│   │   ├── usecase/          # Application logic
│   │   ├── handler/          # HTTP handlers
│   │   ├── router/           # Routing
│   │   └── logging/          # slog helpers / request-ID middleware
│   └── test/                 # Test helpers and mocks
├── db/migrations/            # SQL migrations
├── docs/                     # GitHub Pages content
├── .devcontainer/            # Dev Container definition
├── docker-compose.yml        # api / db / test-db services
├── Dockerfile
└── Makefile
```

## Getting Started

### Prerequisites

- Docker / Docker Compose
- VS Code with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) (recommended)

### 1. Configure environment variables

Copy `.env.example` to `.env`. The defaults work out of the box.

```bash
cp .env.example .env
```

### 2. Start the dev environment (recommended: Dev Container)

Open the repository in VS Code and run **Dev Containers: Reopen in Container** from the command palette. The `api`, `db`, and `test-db` containers start, and your editor attaches to the dev container.

> Without Dev Containers, run `docker compose up -d` and then `docker compose exec api bash` to get a shell inside the container.

### 3. Apply migrations

```bash
make migrate-up
```

### 4. Run the application

```bash
make dev
```

`air` watches the source tree and reloads on changes. Hit `http://localhost:8080/` to get a random user back.

```bash
curl http://localhost:8080/
```

## Make Targets

Run `make` (or `make help`) to see the full list.

| Target | Description |
| --- | --- |
| `make dev` | Run the app with hot reload via `air` |
| `make fmt` | `go fmt ./...` |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down [N=2]` | Roll back migrations (defaults to 1 step) |
| `make migrate-create NAME=add_xxx_table` | Scaffold a new migration pair |
| `make test-setup` | Apply migrations to the test database |
| `make test-repository` | Run repository-layer tests |
| `make test-usecase` | Run usecase-layer tests |
| `make test-handler` | Run handler-layer tests |
| `make test-all` | Run all layered tests |

## Testing

Tests run against the `test-db` service (host port `5433`). Apply migrations once, then run the suite:

```bash
make test-setup
make test-all
```

Use `make test-repository`, `make test-usecase`, or `make test-handler` to scope a run to a single layer. Mocks for each layer live under `app/test/mock/` and are consumed by the handler and usecase tests.

## Adding a New Endpoint

Follow the Clean Architecture layers when adding a feature:

1. Define the domain model in `app/internal/entity/`.
2. Add a migration in `db/migrations/` and run `make migrate-up`.
3. Implement persistence in `app/internal/repository/`.
4. Implement the use case in `app/internal/usecase/`.
5. Implement the HTTP handler in `app/internal/handler/`.
6. Register the route in `app/internal/router/router.go`.
7. Wire the dependencies in `app/cmd/main.go`.

Add tests and mocks alongside each layer.
