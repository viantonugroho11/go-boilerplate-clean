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
7. [README.md](../README.md) — update when the change affects how the project is run, configured, or explored (see § README)
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
| CRUD + Kafka on **Create/Update** (event DTO + log on publish failure) | This file, § Event publishing | `internal/usecase/users`, `events.go` | `/users` |
| State machine + Kafka on save (stricter: return publish error) | [statemachine.md](./statemachine.md) | `internal/usecase/sample` | `POST/PUT /samples` |

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
- Kafka publish on Create/Update: {yes — default for persisted writes|no — document why}
- Status / workflow: {yes|no}
- Money / amount fields: {yes|no}        # if yes, use shopspring/decimal (see § Money)

## Rules
1. Schema from `database/*.sql`; HTTP from `database/openapi.yaml` (do not guess)
2. Usecase → repository/publisher interfaces only; pass `tx *gorm.DB` on repo calls inside transactions (reads and writes)
3. Entity: internal/entity/{domain}/ (no GORM)
4. GORM model: internal/repository/{domain}/model/ (match SQL types)
5. HTTP: handler + dto aligned with `openapi.yaml` `operationId` and schemas
6. Wire: internal/bootstrap/wire.go
7. context.Context as first param on all methods
8. Repository: **every** method `(ctx, tx *gorm.DB, ...)` — reads included; usecase passes `tx` from `Begin`
9. Business errors: `internal/shared/apperrors` + `internal/shared/response`
10. Money/amount: `github.com/shopspring/decimal` — map from SQL `NUMERIC` (see § Money)
11. **Create and Update** must call a usecase `*EventPublisher` after successful persist (see § Event publishing); wire producer in `bootstrap/wire.go`
12. go build ./... must pass
13. Update [README.md](../README.md) so it matches the work shipped (endpoints, config, layout, run commands) — see § README

## Files (adjust to feature)
- [ ] internal/entity/{domain}/{entity}.go
- [ ] internal/repository/{domain}/{domain}_repository.go
- [ ] internal/repository/{domain}/model/{entity}.go
- [ ] internal/repository/{domain}/postgres/repository.go
- [ ] internal/usecase/{domain}/{domain}_usecase.go
- [ ] internal/usecase/{domain}/events.go          # {Domain}EventPublisher interface (required if Create/Update persist data)
- [ ] internal/transport/apis/dto/...        # fields from openapi requestBody/schemas
- [ ] internal/transport/apis/handler/...  # one handler per openapi operationId
- [ ] internal/transport/apis/router.go    # paths/methods from openapi paths
- [ ] internal/bootstrap/wire.go
- [ ] internal/infrastructure/database/postgres/connection.go (Migrate)
- [ ] Kafka: event DTO or publisher (see §5)
- [ ] State machine: states/, saver, on_* (see statemachine.md)
- [ ] [README.md](../README.md) — folder structure, endpoints, env vars, Kafka/consumer, docs links
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

**Dependencies:** `transport → usecase → repository (interface)`; both usecase and repository import `entity`. Infrastructure is wired in `bootstrap`, not imported by usecase. Every repository method accepts `*gorm.DB`; usecase passes `tx` from `begin.BeginRepository`.

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

### Parameter `tx` + helper `dbOrTx` (golden reference: sample)

**Every repository interface method must accept `tx *gorm.DB`** (reads and writes). In the Postgres implementation, use the **`dbOrTx`** helper:

| `tx` | Behavior |
|------|----------|
| `tx != nil` | Run queries on the active transaction (read/write inside `Begin` … `Commit`) |
| `tx == nil` | Fall back to `r.db` (reads outside a transaction, e.g. `UserService.GetByID` → `repo.GetByID(ctx, nil, id)`) |

**Writes** (`Create`, `Add`, `Update`, `Delete`): still **require** `tx != nil` in the implementation (do not fall back to `r.db` for mutations).

```go
func (r *sampleRepository) dbOrTx(tx *gorm.DB) *gorm.DB {
    if tx != nil {
        return tx
    }
    return r.db
}
```

**Example interface + read (`internal/repository/sample/`):**

```go
type SampleRepository interface {
    Add(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
    GetByID(ctx context.Context, tx *gorm.DB, id string) (*entitysample.Sample, error)
    Update(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
}
```

```go
func (r *sampleRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entitysample.Sample, error) {
    var m model.Sample
    err := r.dbOrTx(tx).WithContext(ctx).Where("id = ?", id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, apperrors.ErrSampleNotFound
    }
    if err != nil {
        return nil, err
    }
    s := model.ToEntity(&m)
    return &s, nil
}
```

- Models with `gorm.DeletedAt`: GORM `First` / `Find` already exclude soft-deleted rows; add `AND deleted_at IS NULL` only for raw SQL or when not using GORM soft delete.
- Return domain entities, not GORM models.

**Users:** same `(ctx, tx, …)` interface; implementation can be aligned with `dbOrTx` (today `GetByID`/`List` manually fall back to `r.db` when `tx == nil`).

### Interface — `internal/repository/{domain}/{domain}_repository.go`

- Define the interface here; postgres is the only implementation.
- Every method includes `tx *gorm.DB`.

### Model — `internal/repository/{domain}/model/`

- GORM tags, `TableName()`, `ToEntity`, `ToModel`.

### Postgres — `internal/repository/{domain}/postgres/repository.go`

- `NewXxxRepository(db *gorm.DB)` — store `r.db` for `dbOrTx` fallback.
- `dbOrTx(tx *gorm.DB) *gorm.DB` on every postgres repository struct.
- Read: `r.dbOrTx(tx).WithContext(ctx).…`
- Write: `tx.WithContext(ctx).…` with `tx == nil` → error.
- Create: generate UUID when ID is empty (per domain rules).
- `gorm.ErrRecordNotFound` → `apperrors` domain.

### Migrate

Add model to `Migrate()` in `internal/infrastructure/database/postgres/connection.go`:

```go
return db.AutoMigrate(&model.User{}, &samplemodel.Sample{})
```

### Transactions — `internal/repository/begin/`

Usecase **must** open a transaction before calling the repository for coordinated CRUD or state-machine flows:

```go
tx, err := s.txManager.Begin(ctx)
if err != nil { return ..., err }
defer func() {
    if err != nil {
        _ = s.txManager.Rollback(ctx, tx)
    }
}()

// Inside a transaction: pass tx (reads share the same snapshot):
current, err := s.getter.Get(ctx, tx, id)   // → repo.GetByID(ctx, tx, id) → dbOrTx(tx)
updated, err := s.repo.Update(ctx, tx, entity)
// ...

err = s.txManager.Commit(ctx, tx)
```

`begin.BeginRepository` satisfies sample’s `TransactionManager` interface — wire it directly in `wireSampleService`.

**Outside a transaction:** pass `nil` for reads — `repo.GetByID(ctx, nil, id)` uses `r.db` via `dbOrTx`.

**Getter:** `SampleGetter.Get(ctx, tx, id)` forwards `tx` to the repository (see `internal/usecase/sample/getter.go`; `saver.go` calls it after `Begin`).

---

## 3. Usecase (`internal/usecase/{domain}/`)

### Event publishing on Create & Update (required — golden reference: users)

Whenever a usecase **persists** data on **Create** or **Update**, it must publish to Kafka through a **publisher interface** defined in the usecase package (not infra imported from usecase).

| Rule | Detail |
|------|--------|
| Interface file | `internal/usecase/{domain}/events.go` — e.g. `UserEventPublisher` |
| Constructor | Inject publisher: `NewUserService(repo, txManager, publisher)` |
| When to publish | **After** DB success (`Commit` for transactional Create; immediately after `repo.Update` / `repo.Create` succeeds) |
| Publish failure | **Log and return success** to HTTP caller (do not roll back DB); see `user_usecase.go` |
| Delete / read-only | No publish unless product spec requires it |
| State machine domains | Still publish on successful save — sample **returns** publish error (stricter); new CRUD domains follow **users** unless `statemachine.md` says otherwise |

**`events.go` (interface only):**

```go
// internal/usecase/users/events.go
type UserEventPublisher interface {
    PublishUser(ctx context.Context, user userEntity.User) error
}
```

**Create — transaction then publish** (`internal/usecase/users/user_usecase.go`):

```go
created, err := s.repo.Create(ctx, tx, user)
if err != nil {
    return userEntity.User{}, err
}
if err = s.txManager.Commit(ctx, tx); err != nil {
    return userEntity.User{}, err
}

if err = s.publisher.PublishUser(ctx, created); err != nil {
    log.Printf("user_usecase: PublishUserCreated: %v", err)
}
return created, nil
```

**Update — persist then publish** (same publisher call pattern):

```go
updated, err := s.repo.Update(ctx, tx, user)
if err != nil {
    return userEntity.User{}, err
}

if err = s.publisher.PublishUser(ctx, updated); err != nil {
    log.Printf("user_usecase: PublishUserUpdated: %v", err)
}
return updated, nil
```

**Wiring checklist per domain:**

1. `internal/transport/event/events/{domain}_events.go` — **one** Kafka JSON DTO per domain (e.g. `CampaignEvent`; legacy users: `UserCreatedEvent`)
2. `internal/infrastructure/broker/kafka/{domain}_event_publisher.go` — implements `{Domain}EventPublisher` using `go-lib/kafka` `Producer[T]`
3. `internal/bootstrap/wire.go` — `NewProducer`, `New{Domain}EventPublisherKafka`, pass into `New{Domain}Service`
4. `internal/usecase/{domain}/{domain}_usecase.go` — call publisher on **Create** and **Update** only

### Event naming: one `{Domain}Event` per domain vs **sample** exception

Use the **same call pattern** as users (publish after persist; interface in `events.go`; wire one `Producer[{Domain}Event]` per domain). Use **one event struct per domain** — not `CampaignCreatedEvent` / `CampaignUpdatedEvent`.

| | **New CRUD domains** (recommended) | **users** (legacy boilerplate) | **sample** (state machine — exception) |
|---|-----------------------------------|----------------------------------|----------------------------------------|
| **When to publish** | After Create & Update succeed | Same | After `Commit` in `Save` |
| **Event DTO** | Single `{Domain}Event` + `event_type` (`created` \| `updated`) | `UserCreatedEvent` (no `event_type`; legacy name) | Entity payload `entitysample.Sample` |
| **Publisher API** | `PublishCampaign(ctx, campaign, eventType)` or `Publish(ctx, events.CampaignEvent)` | `PublishUser(ctx, user)` | `Publish(ctx, sample)` |
| **Producer** | `kafka.NewProducer[events.CampaignEvent](...)` — **one type per domain** | `Producer[events.UserCreatedEvent]` | `Producer[entitysample.Sample]` |
| **Publish fails** | **Log**; HTTP still **2xx** | Same (log only) | **Return error** from `Save` |
| **Reference code** | Implement per this table | `user_usecase.go` | `saver.go` (lines 126–127) |

**Recommended layout (e.g. `campaign`):**

```go
// internal/transport/event/events/campaign_events.go
const (
    CampaignEventTypeCreated = "created"
    CampaignEventTypeUpdated = "updated"
)

type CampaignEvent struct {
    EventType string `json:"event_type"` // created | updated
    ID        string `json:"id"`
    // shared fields consumers need
}

// internal/usecase/campaign/events.go
type CampaignEventPublisher interface {
    Publish(ctx context.Context, c campaignEntity.Campaign, eventType string) error
}
```

```go
// infra: maps entity + eventType → CampaignEvent → Producer.Publish
// internal/infrastructure/broker/kafka/campaign_event_publisher.go
func (p *CampaignEventPublisherKafka) Publish(ctx context.Context, c campaign.Campaign, eventType string) error {
    return p.producer.Publish(ctx, events.CampaignEvent{
        EventType: eventType,
        ID:        c.ID,
        // ...
    })
}
```

```go
// Create — after Commit
if err = s.publisher.Publish(ctx, created, events.CampaignEventTypeCreated); err != nil {
    log.Printf("campaign_usecase: publish created: %v", err)
}

// Update — after repo.Update
if err = s.publisher.Publish(ctx, updated, events.CampaignEventTypeUpdated); err != nil {
    log.Printf("campaign_usecase: publish updated: %v", err)
}
```

One topic per domain is typical; consumers branch on `event_type` in JSON. Document topic name in README / config.

**Do not** use `{Domain}CreatedEvent` / `{Domain}UpdatedEvent` for new domains. **Do not** apply sample’s “return publish error” to new CRUD unless product requires it — see [statemachine.md](./statemachine.md) § Kafka publish semantics.

**users** stays as legacy naming (`UserCreatedEvent`); new domains use `{Domain}Event` + `event_type`.

### CRUD service (users)

- Files: `{domain}_usecase.go` + **`events.go`** (publisher interface).
- `NewUserService(repo, txManager, publisher)` — **must** pass `begin.BeginRepository` and `UserEventPublisher`.
- Mutations: `Begin` → repo writes with `tx`; reads in tx: `repo.GetByID(ctx, tx, id)`; reads without tx: `repo.GetByID(ctx, nil, id)` (falls back to `r.db`).
- **Create/Update:** always `publisher.Publish…` after successful persist; log publish errors (see § Event publishing).

### State machine service (sample) — publish-error exception

- `SampleService` with `Save(ctx, sample)` — see [statemachine.md](./statemachine.md).
- `NewSampleSaver(factory, txManager, adder, getter, updater, publisher)`.
- Inside `Save`: `getter.Get(ctx, tx, id)` → `repo.GetByID(ctx, tx, id)`; `adder`/`updater` use `(ctx, tx, ...)` as well.
- Publish after `Commit`; if `publisher.Publish` fails, **`Save` returns the error** (DB already committed — operators must handle retry/idempotency). This is **intentionally stricter** than CRUD users; do not copy this into new CRUD usecases.

### Publisher interface (in usecase package)

| Domain | Interface location | Method | Infra implementation |
|--------|-------------------|--------|----------------------|
| users | `events.go` | `PublishUser(ctx, user)` | `user_event_publisher.go` → `Producer[events.UserCreatedEvent]` |
| sample | `saver.go` | `Publish(ctx, sample)` | `sample_event_publisher.go` → `Producer[entitysample.Sample]` |

Usecase must depend on the **interface**, never on `internal/infrastructure/broker/kafka` directly.

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

## 5. Kafka producers (required for Create & Update)

**Default:** any domain with HTTP (or workflow) **Create** or **Update** that writes to Postgres must ship a producer wired like **users**.

| Style | Used by | When | Event type | Publisher |
|-------|---------|------|------------|-----------|
| Legacy single DTO | users | After Create & Update | `UserCreatedEvent` (no `event_type`) | `user_event_publisher.go` |
| **One DTO per domain** | **new CRUD** (recommended) | After Create / Update | `{Domain}Event` + `event_type` `created`\|`updated` | `{domain}_event_publisher.go` |
| Domain entity as payload | sample (state machine) | After `Save` commits | `entitysample.Sample`; **return** publish error | `sample_event_publisher.go` |

**Not optional for new CRUD domains** unless the task explicitly says “no Kafka”. If omitted, document the reason in the PR/README.

**Producer flow (users):**

```
usecase Create/Update success
    → publisher.Publish*(ctx, entity)
        → kafka/*_event_publisher.go maps entity → events.*Event
            → go-lib Producer.Publish → topic from config
```

Shared setup in `wire.go` (per domain topic/key as needed):

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

## 7. README.md (keep in sync with the project)

[README.md](../README.md) is the **human-facing** overview: how to run the app, configure env vars, and discover APIs. Agents must **update it in the same task** when codegen changes what operators or new contributors need to know.

**Language:** English (same as this doc).

**Do not** paste full OpenAPI or SQL into README — link to `database/openapi.yaml`, `database/README.md`, and `database/*.sql` instead.

### When README must be updated

| You changed | Update README section(s) |
|-------------|---------------------------|
| New HTTP routes / domains | **HTTP Endpoints** (table + example cURL); mention in intro if it is a primary API |
| `database/openapi.yaml` paths (e.g. campaigns) | **HTTP Endpoints** + **Code generation** bullet; link to `database/openapi.yaml` |
| New `cmd/*` entrypoint | **Folder Structure**, **Run Locally** (commands) |
| Kafka topic / consumer / producer | **Kafka** (library, paths, env vars, how to run consumer) |
| New env keys in `configs/config.yaml` or `internal/config` | **Configuration** |
| New tables / `database/*.sql` | **Database** (schema source: `database/`, AutoMigrate models) |
| New reference domain in repo | **Repository & Usecase** or a short **Domains** list |
| Bootstrap / wiring pattern | **Architecture Notes**, startup behavior in **Run Locally** |
| Stale paths in README (e.g. `transport/http` → `transport/apis`) | Fix while touching README for any reason |

### Verify before editing README

Read the **actual** tree and wiring — do not copy outdated README text:

| Topic | Current boilerplate (verify in repo) |
|-------|--------------------------------------|
| HTTP | `internal/transport/apis/` (`handler/`, `dto/`, `router.go`) |
| Entrypoints | `cmd/app` (HTTP), `cmd/consumer` (Kafka consumer, `-consumer=user\|order`) |
| Bootstrap | `internal/bootstrap/` (`app.go`, `wire.go`, `server.go`, `consumer.go`) |
| Kafka library | `github.com/viantonugroho11/go-lib/kafka` (not raw Sarama in app code) |
| Schema / API contract | `database/*.sql`, `database/openapi.yaml` |
| Example APIs | `/users`, `POST/PUT /samples`; campaigns per `openapi.yaml` when implemented |
| Transactions | `internal/repository/begin/`, repo methods use `tx *gorm.DB` + `dbOrTx` |

### Suggested README sections (template)

Keep existing headings where possible; extend rather than rewrite the whole file.

1. **Title & stack** — Echo, GORM, Postgres, Kafka (go-lib), Redis, Viper; `database/`-first workflow.
2. **Prerequisites** — unchanged unless new infra (e.g. another broker).
3. **Folder Structure** — reflect `cmd/app`, `cmd/consumer`, `internal/bootstrap`, `database/`, `internal/entity`, `internal/repository/{domain}`, `internal/usecase`, `internal/transport/apis`, `internal/transport/event`.
4. **Configuration** — env vars from `configs/config.yaml` / `internal/config` (include new topics e.g. `KAFKA_TOPIC_ORDERS` when added).
5. **Run Locally** — `go run ./cmd/app`; consumers: `go run ./cmd/consumer -consumer=user`.
6. **HTTP Endpoints** — tables: Method, Path, Notes; separate **Boilerplate examples** vs **OpenAPI domains** (campaigns).
7. **Kafka** — producer in `internal/infrastructure/broker/kafka/`; consumers in `internal/transport/event/kafka/`; wire in `bootstrap`.
8. **Database** — SQL in `database/`; GORM AutoMigrate in `connection.go`; optional link to `database/README.md`.
9. **Code generation (AI agents)** — links to `AGENTS.md`, `docs/codegen.md`, `docs/statemachine.md`, `database/`.
10. **Architecture Notes** — layer rules + pointer to `dbOrTx` / transactions (one short paragraph).

### Example: adding domain `campaign`

After implementing from `database/*.sql` + `openapi.yaml`:

```markdown
## HTTP Endpoints

### Campaigns (OpenAPI: `database/openapi.yaml`)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/categories` | List campaign categories |
| GET | `/campaigns` | List campaigns |
| POST | `/campaigns` | Create campaign |
| ... | ... | ... |

Schema and request bodies: see `database/openapi.yaml`.
```

Also ensure **Folder Structure** mentions `internal/entity/campaigns`, `internal/repository/campaign`, etc., if those packages exist.

### README checklist (per feature)

- [ ] Paths and package names match the repo (no `transport/http`, no wrong `main.go` paths)
- [ ] New routes and env vars documented
- [ ] Consumer command documented if a new `-consumer=` name was added
- [ ] Links to `database/openapi.yaml` / `database/README.md` for contract details
- [ ] **Code generation** section still points to `AGENTS.md` and `docs/codegen.md`
- [ ] No secrets or `.env` contents in README

---

## 8. Pre-completion checklist

- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] Routes + `RegisterRoutes` signature updated
- [ ] `wire{Domain}Service` added; `app.go` / `server.go` pass new service
- [ ] `AutoMigrate` includes new model
- [ ] `txManager` injected where transactions are used
- [ ] Every repository method has `tx *gorm.DB`; postgres has `dbOrTx`; writes require `tx != nil`
- [ ] `events.go` + `{domain}_event_publisher.go` + `events/{domain}_events.go`; **Create** and **Update** call publisher; publish errors logged (users style)
- [ ] `wire{Domain}Service` creates `kafka.NewProducer` and injects publisher
- [ ] State machine domains: [statemachine.md](./statemachine.md)
- [ ] Money fields use `shopspring/decimal`, not float
- [ ] Read `database/*.sql` + `database/openapi.yaml`; implementation matches both
- [ ] [README.md](../README.md) updated for endpoints, config, run commands, and layout (§ README)

---

## AI command examples

> Read AGENTS.md and docs/codegen.md. Add **product** with HTTP CRUD and postgres, **including Kafka publish on Create/Update** (users `events.go` + publisher pattern).

> Read AGENTS.md, docs/codegen.md, and docs/statemachine.md. Add **order** with status workflow like **sample**, including Kafka and POST/PUT routes.

> Add **invoice** with fields `amount` and `tax` using shopspring/decimal per docs/codegen.md (§ Money). Postgres numeric(18,2), no float64.

> Implement **campaigns** API: read `database/README.md`, all `database/*.sql`, and **`database/openapi.yaml`**. Generate entity/repo from SQL; handlers/DTOs/routes **only** from openapi `operationId` and schemas. Use `decimal` for `target_amount` / `collected_amount`. Update **README.md** (HTTP endpoints, folder structure, config). Run `go build ./...`.
