# Agent instructions — go-boilerplate-clean

Read this file first. Then open the doc(s) for your task before writing code.

## Documentation map

| Task | Read | Golden reference |
|------|------|------------------|
| Schema / persistence | `database/*.sql`, [database/README.md](database/README.md) | SQL is source of truth |
| HTTP handlers / DTOs / routes | **`database/openapi.yaml`** (required) + [docs/codegen.md](docs/codegen.md) | OpenAPI `operationId` + schemas |
| Go patterns (CRUD, wire) | [docs/codegen.md](docs/codegen.md) | `internal/usecase/users` |
| Entity with `status` / workflow | [docs/statemachine.md](docs/statemachine.md) | `internal/usecase/sample` |

**Module path:** `go-boilerplate-clean`

## Reference implementations (current repo)

### Users — simple CRUD + Kafka

| Layer | Path |
|-------|------|
| Entity | `internal/entity/users/user.go` |
| Repository | `internal/repository/user/` |
| Usecase | `internal/usecase/users/user_usecase.go`, `events.go` |
| HTTP | `internal/transport/apis/handler/user_handler.go` |
| Kafka event | `internal/transport/event/events/user_events.go` |
| Kafka publisher | `internal/infrastructure/broker/kafka/user_event_publisher.go` |
| Wire | `wireUserService` in `internal/bootstrap/wire.go` |

### Sample — state machine + Kafka

| Layer | Path |
|-------|------|
| Entity + status constants | `internal/entity/sample/sample.go` |
| Repository | `internal/repository/sample/` (`Add`, `GetByID`, `Update`) |
| Orchestration | `internal/usecase/sample/saver.go` (`SampleService.Save`) |
| State machine | `internal/usecase/sample/states/` |
| Transition handlers | `on_open/`, `on_pending/`, `on_closed/` |
| HTTP | `POST /samples`, `PUT /samples/:id` |
| Kafka publisher | `internal/infrastructure/broker/kafka/sample_event_publisher.go` (payload: entity) |
| Wire | `wireSampleService` in `internal/bootstrap/wire.go` |

## Required workflow

1. Read **`database/README.md`**, relevant **`database/*.sql`**, and for any HTTP work **`database/openapi.yaml`** (mandatory for handlers/DTOs/routes).
2. Read [docs/codegen.md](docs/codegen.md) and/or [docs/statemachine.md](docs/statemachine.md).
3. Respect layer boundaries: `entity` → `repository` → `usecase` → `transport` → `bootstrap`.
4. Usecase depends on **interfaces** only (repository, publisher, `begin.BeginRepository` / `TransactionManager`).
5. Register routes in `internal/transport/apis/router.go` — paths/methods from **`database/openapi.yaml`**.
6. Wire in `internal/bootstrap/wire.go` (`wire{Domain}Service`).
7. Register GORM models in `internal/infrastructure/database/postgres/connection.go` → `Migrate()` (schema from **`database/*.sql`**).
8. Finish with `go build ./...` (and `go vet ./...` when possible).

## Conventions

- **Entity package** may be plural (`users`); **repository folder** is often singular (`user`, `sample`).
- **Transactions:** inject `begin.BeginRepository` — it matches sample’s `TransactionManager` (`Begin` / `Commit` / `Rollback`).
- **Kafka — publish on Create/Update (CRUD) or after Save (state machine):** see [docs/codegen.md](docs/codegen.md) § Event publishing.
  - **New CRUD domains:** one `{Domain}Event` per domain with `event_type` (`created` / `updated`); same call pattern as users; **log** publish failures (HTTP still succeeds).
  - **Users (legacy):** `UserCreatedEvent` (legacy name, no `event_type`) — call pattern only.
  - **Sample (exception):** entity payload; `Save` **returns** publish error to caller ([statemachine.md](docs/statemachine.md) § Kafka publish semantics).
- **State machine:** all status changes go through `states` + `Saver.Save`; never in HTTP handlers or repository shortcuts.

## HTTP endpoints (existing)

| Method | Path | Service |
|--------|------|---------|
| CRUD | `/users` | `UserService` |
| Create / update workflow | `POST /samples`, `PUT /samples/:id` | `SampleService` |
| Health | `GET /healthz` | — |

## Do not

- Put GORM tags or DB logic in `internal/entity/*` (JSON tags on entity are OK when returned as API JSON, e.g. sample).
- Put status transition logic in Echo handlers.
- Import `postgres` repository implementations from `usecase` or `states` packages.
- Commit `.env` or secrets.

## Example agent prompts

```text
Read AGENTS.md and docs/codegen.md. Add domain "product" with HTTP CRUD and postgres repo. No Kafka. Run go build ./...
```

```text
Read AGENTS.md, docs/codegen.md, and docs/statemachine.md. Add domain "order" with status workflow like sample. Wire Kafka. Reference internal/usecase/sample. Run go build ./...
```
