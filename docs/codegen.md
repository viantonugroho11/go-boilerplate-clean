# Codegen Prompt — go-boilerplate-clean

This document is a **ready-to-use prompt** for AI/codegen when adding new features to this boilerplate. Follow existing repo conventions; do not change the architecture without a strong reason.

**Module path:** `go-boilerplate-clean`

**Reference implementations:**
- Simple CRUD: `internal/usecase/users`, `internal/repository/user`, `internal/transport/apis`
- State machine (full stack): `internal/usecase/sample`, `internal/repository/sample`, `POST/PUT /samples` → [statemachine.md](./statemachine.md)

**Agent entry point:** [AGENTS.md](../AGENTS.md)

---

## Prompt template (copy & fill placeholders)

```markdown
You are adding a new feature to the Go repo `go-boilerplate-clean` (clean architecture).

## Feature context
- Domain name (snake_case): {domain}          # e.g. order, product
- Entity name (PascalCase): {Entity}           # e.g. Order, Product
- Operations: {create|read|update|delete|list|custom}
- Needs DB transaction: {yes|no}
- Needs Kafka/event publish: {yes|no}
- Has status / workflow field: {yes|no}        # if yes, read docs/statemachine.md

## Required rules
1. Follow folder structure and patterns in docs/codegen.md
2. Usecase depends only on repository/publisher **interfaces**, not postgres/kafka implementations
3. Domain entity in `internal/entity/{domain}/`, no GORM tags
4. GORM model in `internal/repository/{domain}/model/`
5. Echo handler in `internal/transport/apis/handler/`, DTO in `internal/transport/apis/dto/`
6. Wire dependencies in `internal/bootstrap/wire.go` (or a separate wire function per domain)
7. Every method takes `context.Context` as the first argument
8. Business errors: `errors.New("clear message")` — do not expose DB details over HTTP
9. Do not commit `.env` files or secrets

## Expected output
Create/update these files (adjust the list to the feature):
- [ ] internal/entity/{domain}/{entity}.go
- [ ] internal/repository/{domain}/{domain}_repository.go
- [ ] internal/repository/{domain}/model/{entity}.go
- [ ] internal/repository/{domain}/postgres/repository.go
- [ ] internal/usecase/{domain}/{domain}_usecase.go (+ events.go if publishing)
- [ ] internal/transport/apis/dto/{domain}_request.go
- [ ] internal/transport/apis/handler/{domain}_handler.go
- [ ] internal/transport/apis/router.go (register routes)
- [ ] internal/bootstrap/wire.go (inject repo, publisher, usecase)
- [ ] internal/transport/event/events/{domain}_events.go (if Kafka)
- [ ] internal/infrastructure/broker/kafka/{domain}_event_publisher.go (if Kafka)
- [ ] internal/usecase/{domain}/states/* + saver (if state machine)

After generation, ensure `go build ./...` passes.
```

---

## Layer architecture

```
cmd/app/main.go
    └── internal/
            ├── bootstrap/      # config, DB, redis, wire, HTTP server
            ├── config/         # configuration loader (Viper / go-config-library)
            ├── entity/         # domain model (pure, no GORM)
            ├── usecase/        # business logic, validation, tx & event orchestration
            ├── repository/     # interface + postgres impl + model mapping
            ├── transport/      # HTTP (Echo), event (Kafka consumer)
            └── infrastructure/ # DB, broker, cache connections
```

**Dependency direction:** `transport → usecase → repository (interface)`; `usecase` and `repository` import `entity`. Infrastructure is injected from bootstrap, not imported directly by usecase (except `*gorm.DB` via transaction parameters).

---

## Naming conventions

| Item | Pattern | Example |
|------|---------|---------|
| Domain folder | lowercase singular (entity may be plural) | entity: `users`, repo: `user`, `sample` |
| Entity package | same as folder | `package users` |
| Entity struct | PascalCase | `User`, `Sample` |
| Repository interface | `{Domain}Repository` | `UserRepository` |
| Postgres implementation | `internal/repository/{domain}/postgres` | `NewUserRepository(db)` |
| Usecase service | `{Domain}Service` + `New{Domain}Service` | `UserService` |
| HTTP handler | `{Domain}Handler` | `UserHandler` |
| Request DTO | `{Action}{Domain}Request` | `CreateUserRequest` |
| Event payload | `{Domain}{Action}Event` | `UserCreatedEvent` |
| Status constant | `{Entity}Status{Name}` | `SampleStatusOpen` |

Entity import alias (when package names collide):

```go
userEntity "go-boilerplate-clean/internal/entity/users"
```

---

## 1. Entity (`internal/entity/{domain}/`)

- Domain struct **without** GORM / JSON tags unless the entity is returned directly in API responses.
- Status/workflow constants live in the entity (see `internal/entity/sample/sample.go`).
- No DB logic in the entity.

```go
package sample

type Sample struct {
    ID     string
    Status string
    // ...
}

const (
    SampleStatusOpen   = "open"
    SampleStatusOnHold = "on_hold"
    SampleStatusClosed = "closed"
)
```

---

## 2. Repository

### Interface — `internal/repository/{domain}/{domain}_repository.go`

- Every method: `(ctx context.Context, tx *gorm.DB, ...)`.
- `tx` may be `nil` for reads outside a transaction (see `GetByID` / `List` in the user repo).
- Return type: domain entity, not the GORM model.

### Model — `internal/repository/{domain}/model/`

- Struct with GORM tags as needed.
- Explicit `TableName()`.
- `ToEntity(*Model) Entity` and `ToModel(Entity) Model`.

### Postgres — `internal/repository/{domain}/postgres/repository.go`

- Private struct `xxxRepository` with field `db *gorm.DB`.
- Constructor: `NewXxxRepository(db *gorm.DB) xxx.XxxRepository`.
- Create: generate UUID when ID is empty (`github.com/google/uuid`).
- Update/Delete inside a transaction: require `tx != nil`, return an error otherwise.
- `gorm.ErrRecordNotFound` → business error `"xxx not found"`.

### Transactions — `internal/repository/begin/`

- Use `begin.BeginRepository` for `Begin` / `Commit` / `Rollback`.
- Usecase pattern (users):

```go
tx, err := s.txManager.Begin(ctx)
if err != nil { return ..., err }
defer func() {
    if err != nil {
        _ = s.txManager.Rollback(ctx, tx)
    }
}()
// ... repo operations with tx ...
err = s.txManager.Commit(ctx, tx)
```

---

## 3. Usecase (`internal/usecase/{domain}/`)

### CRUD service (no state machine)

Example: `internal/usecase/users/user_usecase.go`

- Public interface: `{Domain}Service` with business operations.
- Private struct: `{domain}Service` — fields: repo, `txManager`, publisher (optional).
- Constructor must inject `begin.BeginRepository` as `txManager` when using transactions (see `NewUserService`).
- Validate input in the usecase (not in the handler).
- After a successful commit, publish events; on publish failure **log only**, do not roll back the transaction (unless requirements say otherwise).

### Publisher interface — `events.go`

```go
type UserEventPublisher interface {
    PublishUser(ctx context.Context, user userEntity.User) error
}
```

Kafka implementation lives in `internal/infrastructure/broker/kafka/` and implements the usecase interface.

### Usecase with state machine

Do not put status transition logic in the handler. Use:

- `Saver` / orchestration service: `internal/usecase/sample/saver.go`
- State machine: `internal/usecase/{domain}/states/`
- Per-transition handlers: `internal/usecase/{domain}/on_{state}/`

Full details: [statemachine.md](./statemachine.md).

---

## 4. HTTP transport

### DTO — `internal/transport/apis/dto/`

- Request structs with `json` tags.
- `ToEntity()` method to map to the entity (no ID on create).

### Handler — `internal/transport/apis/handler/`

- Inject `{Domain}Service` via constructor.
- `c.Bind(&req)` → validation failure → `400`.
- Call usecase with `c.Request().Context()`.
- HTTP status mapping: create `201`, get/list/update `200`, delete `204`, not found `404`, validation/business `400`.

### Router — `internal/transport/apis/router.go`

- Route group: `e.Group("/{plural}")`.
- Register in `RegisterRoutes(e, services...)`.

---

## 5. Kafka events (optional)

| Layer | Path |
|-------|------|
| Payload | `internal/transport/event/events/{domain}_events.go` |
| Publisher impl | `internal/infrastructure/broker/kafka/{domain}_event_publisher.go` |
| Consumer handler | `internal/transport/event/kafka/` |
| Consumer wiring | `internal/bootstrap/consumer.go`, `cmd/consumer/main.go` |

Producer uses `github.com/viantonugroho11/go-lib/kafka` — see `internal/bootstrap/wire.go`.

---

## 6. Bootstrap wiring

- `internal/bootstrap/app.go` — init order: config → DB → wire services → redis → HTTP.
- Add `wire{Domain}Service(db)` in `wire.go`:
  1. Postgres repo
  2. Kafka producer + publisher (if needed)
  3. `New{Domain}Service(...)`
  4. Return `(Service, cleanup func(), error)`

See `wireSampleService` in `internal/bootstrap/wire.go` for the full sample stack (repo, factory, saver, routes).

---

## 7. Pre-completion checklist

- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] Routes registered & handlers wired
- [ ] Repository interface does not depend on the postgres package
- [ ] Entity does not import GORM
- [ ] Transactions: defer rollback only on error (follow existing pattern)
- [ ] Status/workflow: follow [statemachine.md](./statemachine.md)

---

## Short AI command examples

> "Add domain **product** with HTTP CRUD, postgres repository, no Kafka. Follow docs/codegen.md."

> "Add domain **order** with a **status** field, save via state machine. Follow docs/codegen.md and docs/statemachine.md; reference: internal/usecase/sample."
