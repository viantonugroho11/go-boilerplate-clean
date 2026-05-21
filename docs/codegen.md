# Codegen Prompt — go-boilerplate-clean

Ready-to-use guidance for AI/codegen when adding features. Follow repo conventions; avoid architectural changes without good reason.

**Module path:** `go-boilerplate-clean`  
**Agent entry point:** [AGENTS.md](../AGENTS.md)

## Database-first codegen (required)

Before generating **entity, repository, usecase, or migrations**, read the project **`database/`** folder. Before generating or changing **HTTP handlers, DTOs, or routes**, you **must** read **`database/openapi.yaml`**.

| File / folder | Purpose | Authority |
|---------------|---------|-----------|
| [`database/*.sql`](../database/) | PostgreSQL schema (tables, constraints, indexes) | **Source of truth** for persistence |
| [`database/README.md`](../database/README.md) | Domain overview, entity list, execution order | Start here |
| [`database/dbdiagram.dbml`](../database/dbdiagram.dbml) | ERD / relationships (visual reference) | Diagram only |
| [`database/openapi.yaml`](../database/openapi.yaml) | REST contract (paths, methods, schemas, `operationId`) | **Source of truth** for HTTP layer |

**Conflict resolution:** if `openapi.yaml`, `dbdiagram.dbml`, and `.sql` disagree, follow **`database/*.sql`** for columns/types/constraints and **`database/openapi.yaml`** for API paths, request/response shapes, and status codes.

### Read order for agents

```
1. database/README.md
2. database/*.sql          → entity, GORM model, repository, Migrate
3. database/dbdiagram.dbml → relationships (optional)
4. database/openapi.yaml   → handler, DTO, router (mandatory for HTTP)
5. docs/codegen.md (+ docs/statemachine.md if workflow)
6. Reference Go code (users / sample) for patterns only — do not invent routes
```

### Map SQL → Go (repository layer)

| SQL | Entity field | GORM model |
|-----|--------------|------------|
| `UUID` | `string` | `type:uuid` |
| `VARCHAR(n)` | `string` | `size:n` |
| `TEXT` | `string` | `type:text` |
| `NUMERIC(20,2)` (money) | `decimal.Decimal` | `type:numeric(20,2)` — see § Money |
| `TIMESTAMPTZ` | `time.Time` | `timestamptz` |
| `CHECK (...)` / `ENUM`-like varchar | const or validated string in usecase | same constraint awareness |

Derive `TableName()`, `ToEntity`, `ToModel` from the `.sql` column list. Register models in `internal/infrastructure/database/postgres/connection.go` → `Migrate()`.

### Map OpenAPI → Go (HTTP layer) — mandatory

Read every `paths` entry and `components.schemas` in [`database/openapi.yaml`](../database/openapi.yaml) before writing handlers.

| OpenAPI | Go |
|---------|-----|
| `paths./foo/bar.post` | `POST` route + handler method |
| `operationId` | Handler name hint (`createCampaign` → `CreateCampaign`) |
| `parameters` (path/query) | Echo `c.Param`, `c.QueryParam`, `pagination.ParseQuery` |
| `requestBody.schema` | DTO struct + `json` tags + `ToEntity()` |
| `responses.200/201.schema` | Response type; wrap with `internal/shared/response` |
| `responses.404` | `apperrors.NotFound` via `response.Error` |
| `format: uuid` | `string` + validation in usecase |
| `type: number` on money fields | Prefer `string` in DTO → `decimal.NewFromString` (see § Money) |

**Handlers must not:**

- Invent paths, methods, or payload fields absent from `openapi.yaml`.
- Change URL shapes or HTTP verbs unless `openapi.yaml` is updated first.
- Return ad-hoc JSON maps; use `response.OK`, `response.Created`, `response.Paginated`, `response.Error`.

**Example (campaigns, from current spec):**

| Method | OpenAPI path | `operationId` |
|--------|--------------|---------------|
| `GET` | `/categories` | `listCampaignCategories` |
| `GET` | `/campaigns` | `listCampaigns` |
| `POST` | `/campaigns` | `createCampaign` |
| `GET` | `/campaigns/{id}` | `getCampaign` |
| `PATCH` | `/campaigns/{id}` | `patchCampaign` |
| `POST` | `/campaigns/{id}/updates` | `createCampaignUpdate` |
| `GET` | `/campaigns/{id}/statistics` | `getCampaignStatistics` |

Register routes under the same path prefixes in `internal/transport/apis/router.go` (adjust group prefix if the app mounts `/v1` in bootstrap).

---

## Reference implementations

| Pattern | Docs | Code | HTTP |
|---------|------|------|------|
| CRUD + Kafka event DTO | This file | `internal/usecase/users` | `/users` |
| State machine + Kafka entity payload | [statemachine.md](./statemachine.md) | `internal/usecase/sample` | `POST/PUT /samples` |

---

## Prompt template (copy & fill placeholders)

```markdown
You are adding a feature to Go repo `go-boilerplate-clean` (clean architecture).
Read AGENTS.md and docs/codegen.md first; if status/workflow, also docs/statemachine.md.

## Required inputs (read before coding)
1. `database/README.md` and relevant `database/*.sql` (schema source of truth)
2. `database/openapi.yaml` — **mandatory** for handlers, DTOs, routes
3. `database/dbdiagram.dbml` (optional ERD)

## Feature context
- Domain (snake_case): {domain}
- Entity (PascalCase): {Entity}
- Operations: {create|read|update|delete|list|custom}
- DB transaction: {yes|no}
- Kafka publish: {yes|no}
- Status / workflow: {yes|no}
- Money / amount fields: {yes|no}        # if yes, use shopspring/decimal (see § Money)

## Rules
1. Schema from `database/*.sql`; HTTP from `database/openapi.yaml` (do not guess)
2. Usecase → repository/publisher interfaces only
3. Entity: internal/entity/{domain}/ (no GORM)
4. GORM model: internal/repository/{domain}/model/ (match SQL types)
5. HTTP: handler + dto aligned with `openapi.yaml` `operationId` and schemas
6. Wire: internal/bootstrap/wire.go
7. context.Context as first param on all methods
8. Business errors: `internal/shared/apperrors` + `internal/shared/response`
9. Money/amount: `github.com/shopspring/decimal` — map from SQL `NUMERIC` (see § Money)
10. go build ./... must pass

## Files (adjust to feature)
- [ ] internal/entity/{domain}/{entity}.go
- [ ] internal/repository/{domain}/{domain}_repository.go
- [ ] internal/repository/{domain}/model/{entity}.go
- [ ] internal/repository/{domain}/postgres/repository.go
- [ ] internal/usecase/{domain}/...
- [ ] internal/transport/apis/dto/...        # fields from openapi requestBody/schemas
- [ ] internal/transport/apis/handler/...  # one handler per openapi operationId
- [ ] internal/transport/apis/router.go    # paths/methods from openapi paths
- [ ] internal/bootstrap/wire.go
- [ ] internal/infrastructure/database/postgres/connection.go (Migrate)
- [ ] Kafka: event DTO or publisher (see §5)
- [ ] State machine: states/, saver, on_* (see statemachine.md)
```

---

## Layer architecture

```
cmd/app/main.go
    └── internal/
            ├── bootstrap/      # config, DB, wire, HTTP, consumer
            ├── config/
            ├── entity/         # domain model (no GORM)
            ├── usecase/
            ├── repository/     # interface + model + postgres
            ├── transport/      # apis (Echo), event (Kafka)
            └── infrastructure/ # postgres, kafka, redis
```

**Dependencies:** `transport → usecase → repository (interface)`; both usecase and repository import `entity`. Infrastructure is wired in `bootstrap`, not imported by usecase (except `*gorm.DB` on transaction methods).

---

## Naming conventions

| Item | Pattern | Example |
|------|---------|---------|
| Entity folder | often plural | `users` |
| Repository folder | often singular | `user`, `sample` |
| Repository interface | `{Domain}Repository` | `UserRepository`, `SampleRepository` |
| Usecase service | `{Domain}Service` | `UserService`, `SampleService` |
| Wire function | `wire{Domain}Service` | `wireUserService`, `wireSampleService` |
| HTTP handler | `{Domain}Handler` | `UserHandler`, `SampleHandler` |
| Request DTO | `{Action}{Domain}Request` or `Save{Domain}Request` | `CreateUserRequest`, `SaveSampleRequest` |
| Status constant | `{Entity}Status{Name}` | `SampleStatusOpen` |

```go
userEntity "go-boilerplate-clean/internal/entity/users"
entitysample "go-boilerplate-clean/internal/entity/sample"
```

---

## 1. Entity (`internal/entity/{domain}/`)

- No GORM imports or DB logic.
- Status/workflow constants live here for state-machine domains.
- JSON tags are optional; `sample` uses them because handlers return the entity as JSON.

---

## Money, amount, and decimal

For **money**, **price**, **amount**, **fee**, **tax**, **balance**, **quantity with fractional units**, or any value that must not lose precision, use **[shopspring/decimal](https://github.com/shopspring/decimal)** — not `float64`.

| Do | Don't |
|----|--------|
| `decimal.Decimal` in entity & usecase | `float64` / `float32` for money |
| `decimal.NewFromString("136.02")` for parsing | Binary float literals for currency |
| `NUMERIC(p,s)` / `DECIMAL` in Postgres (GORM) | `float` column types in DB |
| String in JSON API when needed (`"19.99"`) or decimal JSON support | Rely on float JSON encoding |

**Install:** `go get github.com/shopspring/decimal`

### Entity (`internal/entity/{domain}/`)

```go
import "github.com/shopspring/decimal"

type Order struct {
    ID     string
    Amount decimal.Decimal // price, total, fee, etc.
}
```

- Arithmetic: `amount.Mul(qty)`, `amount.Add(tax)`, `amount.Sub(discount)` — always produces new values (immutable).
- Comparisons: `amount.Equal(other)`, `amount.GreaterThan(...)`.
- Parsing user input in usecase/DTO: `decimal.NewFromString(s)` and return `apperrors.Validation(...)` on failure.

### GORM model (`internal/repository/{domain}/model/`)

Match precision/scale from `database/*.sql` (e.g. campaigns use `NUMERIC(20,2)`):

```go
import "github.com/shopspring/decimal"

type Campaign struct {
    TargetAmount    decimal.Decimal `gorm:"type:numeric(20,2);not null"`
    CollectedAmount decimal.Decimal `gorm:"type:numeric(20,2);not null"`
}
```

`shopspring/decimal` implements `sql.Scanner` and `driver.Valuer`, so GORM can persist it directly.

### DTO / HTTP

Prefer string amounts in requests to avoid float JSON issues:

```go
type CreateOrderRequest struct {
    Amount string `json:"amount"` // "136.02"
}

func (r *CreateOrderRequest) AmountDecimal() (decimal.Decimal, error) {
    return decimal.NewFromString(r.Amount)
}
```

Or bind as `decimal.Decimal` if the API documents numeric JSON and you control clients.

### Usecase

- Perform all money math with `decimal.Decimal`.
- Round explicitly when needed (`amount.Round(2)`) — do not assume float rounding.
- When splitting amounts (e.g. divide by 3), allocate remainder explicitly (documented in decimal FAQ).

### Checklist (money fields)

- [ ] `github.com/shopspring/decimal` in `go.mod`
- [ ] Entity & model use `decimal.Decimal`
- [ ] Postgres column `numeric` / `decimal`, not `float`
- [ ] DTO parsing validates format and returns validation errors
- [ ] No `float64` in Kafka events for money (use string or decimal serialization)

---

## 2. Repository

### Interface — `internal/repository/{domain}/{domain}_repository.go`

- Methods use `(ctx context.Context, tx *gorm.DB, ...)` when writing inside a transaction.
- Reads outside a tx may omit `tx` (see `UserRepository.GetByID`, `SampleRepository.GetByID`).
- Return domain entities, not GORM models.

**Sample convention:** `Add`, `GetByID`, `Update` (repository implements `SampleAdder` / `SampleUpdater`; getter wraps `GetByID`).

### Model — `internal/repository/{domain}/model/`

- GORM tags, `TableName()`, `ToEntity`, `ToModel`.

### Postgres — `internal/repository/{domain}/postgres/repository.go`

- `NewXxxRepository(db *gorm.DB)`
- Create: UUID if ID empty
- Writes: require `tx != nil`
- `gorm.ErrRecordNotFound` → `"xxx not found"`

### Migrate

Add model to `Migrate()` in `internal/infrastructure/database/postgres/connection.go`:

```go
return db.AutoMigrate(&model.User{}, &samplemodel.Sample{})
```

### Transactions — `internal/repository/begin/`

```go
tx, err := s.txManager.Begin(ctx)
if err != nil { return ..., err }
defer func() {
    if err != nil {
        _ = s.txManager.Rollback(ctx, tx)
    }
}()
// repo calls with tx
err = s.txManager.Commit(ctx, tx)
```

`begin.BeginRepository` satisfies sample’s `TransactionManager` interface — wire it directly in `wireSampleService`.

---

## 3. Usecase (`internal/usecase/{domain}/`)

### CRUD service (users)

- File: `{domain}_usecase.go`, optional `events.go` for publisher interface.
- `NewUserService(repo, txManager, publisher)` — **must** pass `begin.BeginRepository`.
- Validate in usecase; publish after commit; **log** publish errors (do not fail the HTTP response).

### State machine service (sample)

- `SampleService` with `Save(ctx, sample)` — see [statemachine.md](./statemachine.md).
- `NewSampleSaver(factory, txManager, adder, getter, updater, publisher)`.
- Publish after commit; **return** publish error to caller (stricter than users).

### Publisher interface (in usecase package)

```go
// users/events.go
type UserEventPublisher interface {
    PublishUser(ctx context.Context, user userEntity.User) error
}

// sample — interface in saver.go
type SamplePublisher interface {
    Publish(ctx context.Context, sample entitysample.Sample) error
}
```

Implementations: `internal/infrastructure/broker/kafka/`.

---

## 4. HTTP transport

**Contract:** [`database/openapi.yaml`](../database/openapi.yaml) — read fully before implementing this layer.

### DTO — `internal/transport/apis/dto/`

- Struct names aligned with OpenAPI schema names (e.g. `CreateCampaignRequest`, `PatchCampaignRequest`).
- `json` tag names **must match** OpenAPI property names (`snake_case` in spec).
- `ToEntity()` maps into `internal/entity/{domain}/` (no ID on create unless spec includes it).
- Query parameters from OpenAPI `parameters` → parsed in handler or dedicated query struct.
- Use `internal/shared/response` for all responses (see below).

### Handler — `internal/transport/apis/handler/`

- **One handler method per `operationId`** in `openapi.yaml`.
- Inject usecase/service via constructor; no business logic beyond bind/validate/call service/map response.
- `c.Bind` failures → `response.BindError(c, err)`.
- Service errors → `response.Error(c, err)` (`apperrors.FromError`).
- Success → `response.OK`, `response.Created`, `response.Paginated`, or `response.NoContent` per OpenAPI response code.
- HTTP status codes and response bodies must match `openapi.yaml` (`200`, `201`, `404`, etc.).

### Router — `internal/transport/apis/router.go`

- Register **exact** path and HTTP method from `openapi.yaml` `paths` (e.g. `PATCH /campaigns/:id`, not `PUT` unless spec says so).
- Path parameter names must match OpenAPI (`{id}` → `c.Param("id")`).
- Pass wired services from `bootstrap/server.go`.

```go
func RegisterRoutes(e *echo.Echo, userService users.UserService, sampleService usecasesample.SampleService)
```

Legacy routes (`/users`, `/samples`) exist as boilerplate examples; **new domain endpoints come from `openapi.yaml`**.

### OpenAPI ↔ response envelope

Boilerplate JSON shape (`internal/shared/response`):

```json
{ "success": true, "data": { }, "meta": { } }
{ "success": false, "error": { "code": "NOT_FOUND", "message": "..." } }
```

If `openapi.yaml` documents a raw array or schema without wrapper, either update the spec to match the envelope or document a one-off exception in the handler comment — **default is the shared envelope**.

---

## 5. Kafka (optional)

| Style | Used by | Event type | Publisher |
|-------|---------|------------|-----------|
| Dedicated event DTO | users | `events.UserCreatedEvent` | `user_event_publisher.go` |
| Domain entity as payload | sample | `entitysample.Sample` | `sample_event_publisher.go` |

Shared setup in `wire.go`:

```go
producer, err := kafka.NewProducer[T](
    cfg.KafkaBrokersList(),
    cfg.Kafka.Topic,
    kafka.WithKeyFunc[T](func(e T) []byte { return []byte(e.ID) }),
    kafka.WithIdempotent(),
    kafka.WithRetryMax(5),
)
```

| Layer | Path |
|-------|------|
| Event DTO (if used) | `internal/transport/event/events/{domain}_events.go` |
| Publisher | `internal/infrastructure/broker/kafka/{domain}_event_publisher.go` |
| Consumer | `internal/transport/event/kafka/`, `internal/bootstrap/consumer.go` |

Producer library: `github.com/viantonugroho11/go-lib/kafka`.

---

## 6. Bootstrap wiring

**App** (`internal/bootstrap/app.go`): config → DB → `wireUserService` → `wireSampleService` → redis → HTTP.

**Users** — `wireUserService(db)` returns `(UserService, cleanup, error)`; closes Kafka producer on shutdown.

**Sample** — `wireSampleService(db)`:

1. `samplepg.NewSampleRepository(db)`
2. `beginpg.NewBeginRepository(db)` as `TransactionManager`
3. `states.NewSampleStateMachineFactory(on_open, on_pending, on_closed)`
4. Kafka producer + `NewSampleEventPublisherKafka`
5. `NewSampleSaver(...)` → `SampleService`

Also update `internal/bootstrap/consumer.go` if the new domain needs consumer wiring.

---

## 7. Pre-completion checklist

- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] Routes + `RegisterRoutes` signature updated
- [ ] `wire{Domain}Service` added; `app.go` / `server.go` pass new service
- [ ] `AutoMigrate` includes new model
- [ ] `txManager` injected where transactions are used
- [ ] State machine domains: [statemachine.md](./statemachine.md)
- [ ] Money fields use `shopspring/decimal`, not float
- [ ] Read `database/*.sql` + `database/openapi.yaml`; implementation matches both

---

## AI command examples

> Read AGENTS.md and docs/codegen.md. Add **product** with HTTP CRUD and postgres, no Kafka.

> Read AGENTS.md, docs/codegen.md, and docs/statemachine.md. Add **order** with status workflow like **sample**, including Kafka and POST/PUT routes.

> Add **invoice** with fields `amount` and `tax` using shopspring/decimal per docs/codegen.md (§ Money). Postgres numeric(18,2), no float64.

> Implement **campaigns** API: read `database/README.md`, all `database/*.sql`, and **`database/openapi.yaml`**. Generate entity/repo from SQL; handlers/DTOs/routes **only** from openapi `operationId` and schemas. Use `decimal` for `target_amount` / `collected_amount`. Run `go build ./...`.
