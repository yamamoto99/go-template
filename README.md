# go-template

A template repository for building Web APIs in Go with Echo, GORM, and PostgreSQL, following Clean Architecture. It ships with a Dev Container / Docker Compose setup, `golang-migrate` for schema migrations, `air` for hot reload, a layered test suite with mocks, and **schema-first API development** powered by OpenAPI + `oapi-codegen`.

## Tech Stack

- Go 1.26
- [Echo v4](https://echo.labstack.com/) — HTTP framework
- [GORM](https://gorm.io/) + PostgreSQL 17
- [golang-migrate](https://github.com/golang-migrate/migrate) — schema migrations
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) — generates type-safe Echo handlers, models, and embedded spec from OpenAPI 3
- [oapi-codegen/echo-middleware](https://github.com/oapi-codegen/echo-middleware) — runtime request validation against the OpenAPI spec
- [air](https://github.com/air-verse/air) — hot reload
- [testify](https://github.com/stretchr/testify) — assertions and mocks
- Dev Container (VS Code) + Docker Compose

## Project Structure

```
.
├── api/
│   ├── openapi.yaml          # OpenAPI 3 spec (single source of truth for the API)
│   └── cfg.yaml              # oapi-codegen generation config
├── app/
│   ├── cmd/                  # Entry point (main.go)
│   ├── infrastructure/
│   │   ├── config/           # Environment-based configuration loader
│   │   └── db/               # DB connection setup / teardown
│   ├── internal/
│   │   ├── api/              # Generated API code (types, server interface, route registration)
│   │   ├── entity/           # Domain models
│   │   ├── repository/       # Persistence layer
│   │   ├── usecase/          # Application logic
│   │   ├── handler/          # HTTP handlers (implement the generated StrictServerInterface)
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
| `make gen-api` | Generate API code from `api/openapi.yaml` into `app/internal/api/` |
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

## API Development (Schema-First)

The HTTP API is defined in `api/openapi.yaml` and code is generated from it — **edit the spec first, then regenerate**. Do not hand-edit `app/internal/api/api.gen.go`.

```bash
make gen-api
```

What gets generated (into `app/internal/api/`):

- **Models** — Go structs for every schema in the spec (`User`, `Error`, …).
- **`StrictServerInterface`** — typed handler interface: `(ctx, request) → (response, error)`. Handlers implement this and never touch `echo.Context` directly.
- **Response objects** — one type per `(operation, status)` combination (e.g. `GetRandomUser200JSONResponse`). Returning the wrong shape is a compile error.
- **`RegisterHandlers`** — wires every operation to its Echo route.
- **Embedded spec** — the spec is base64+flate compressed into the binary and accessible via `api.GetSwagger()`.

The codegen tool is pinned via the `tool` directive in `go.mod`, so contributors don't need a separate install step — `go tool oapi-codegen` (invoked from `make gen-api`) just works.

### Request validation

`OapiRequestValidator` middleware is registered in `router.NewRouter` and validates every incoming request against the embedded spec before it reaches a handler. Mismatches are rejected automatically:

- Path or method not defined in the spec → `404`
- Body / query / path parameters that violate the schema → `400`

This means handlers can assume their inputs already satisfy the contract — no defensive parsing or "is this field present?" checks at the handler layer.

### Generated code is committed

`app/internal/api/api.gen.go` is checked into git. This keeps diffs reviewable in PRs and removes codegen from the CI critical path. After editing the spec, run `make gen-api` locally and commit the regenerated file alongside your changes.

## Adding a New Endpoint

The flow is schema-first: describe the endpoint in OpenAPI, regenerate, then implement the layers.

1. **Add the path / schemas in `api/openapi.yaml`.** Give it a stable `operationId` — that name appears in the generated interface.
2. **Run `make gen-api`.** New entries appear in `app/internal/api/`: a method on `StrictServerInterface`, request/response objects, and any new model structs.
3. Define the domain model in `app/internal/entity/` (if it differs from the API model).
4. Add a migration in `db/migrations/` and run `make migrate-up`.
5. Implement persistence in `app/internal/repository/`.
6. Implement the use case in `app/internal/usecase/`.
7. Implement the HTTP handler in `app/internal/handler/` — your handler type must satisfy the new method on `api.StrictServerInterface`. The compiler will tell you if the signature or response shape is wrong.
8. Wire the dependencies in `app/cmd/main.go`.

`app/internal/router/router.go` does **not** need per-endpoint edits — `api.RegisterHandlers` reads every operation from the spec. Add tests and mocks alongside each layer.
