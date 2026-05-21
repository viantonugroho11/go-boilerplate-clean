# Codegen Prompt — go-boilerplate-clean

Ready-to-use guidance for AI/codegen when adding features. Follow repo conventions; avoid architectural changes without good reason.

**Module path:** `go-boilerplate-clean`  
**Agent entry point:** [AGENTS.md](../AGENTS.md)

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

## Feature context
- Domain (snake_case): {domain}
- Entity (PascalCase): {Entity}
- Operations: {create|read|update|delete|list|custom}
- DB transaction: {yes|no}
- Kafka publish: {yes|no}
- Status / workflow: {yes|no}
- Money / amount fields: {yes|no}        # if yes, use shopspring/decimal (see § Money)

## Rules
1. Usecase → repository/publisher interfaces only
2. Entity: internal/entity/{domain}/ (no GORM)
3. GORM model: internal/repository/{domain}/model/
4. HTTP: internal/transport/apis/handler/ + dto/
5. Wire: internal/bootstrap/wire.go
6. context.Context as first param on all methods
7. Business errors: `internal/shared/apperrors` (not raw DB errors to HTTP)
8. Money/amount: use `github.com/shopspring/decimal` — never `float32`/`float64`
9. go build ./... must pass

## Files (adjust to feature)
- [ ] internal/entity/{domain}/{entity}.go
- [ ] internal/repository/{domain}/{domain}_repository.go
- [ ] internal/repository/{domain}/model/{entity}.go
- [ ] internal/repository/{domain}/postgres/repository.go
- [ ] internal/usecase/{domain}/...
- [ ] internal/transport/apis/dto/...
- [ ] internal/transport/apis/handler/...
- [ ] internal/transport/apis/router.go
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

```go
import "github.com/shopspring/decimal"

type Order struct {
    Amount decimal.Decimal `gorm:"type:numeric(18,2);not null"`
}
```

`shopspring/decimal` implements `sql.Scanner` and `driver.Valuer`, so GORM can persist it directly. Pick precision/scale to match business rules (e.g. `numeric(18,2)` for currency).

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

### DTO — `internal/transport/apis/dto/`

- `json` tags on requests; `ToEntity()` without ID on create.

### Handler — `internal/transport/apis/handler/`

- Inject service via constructor.
- `c.Bind` → `400` on failure.
- `c.Request().Context()` into usecase.
- Status: `201` create, `200` OK, `204` delete, `404` not found, `400` business/validation.

### Router — `internal/transport/apis/router.go`

```go
func RegisterRoutes(e *echo.Echo, userService users.UserService, sampleService usecasesample.SampleService)
```

Add groups and pass services from `bootstrap/server.go`.

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

---

## AI command examples

> Read AGENTS.md and docs/codegen.md. Add **product** with HTTP CRUD and postgres, no Kafka.

> Read AGENTS.md, docs/codegen.md, and docs/statemachine.md. Add **order** with status workflow like **sample**, including Kafka and POST/PUT routes.

> Add **invoice** with fields `amount` and `tax` using shopspring/decimal per docs/codegen.md (§ Money). Postgres numeric(18,2), no float64.
