# Agent instructions — go-boilerplate-clean

Before generating or modifying features in this repository, read:

| Task | Document | Reference code |
|------|----------|----------------|
| New domain (CRUD, HTTP, repo, wire) | [docs/codegen.md](docs/codegen.md) | `internal/usecase/users` |
| Entity with `status` / workflow | [docs/statemachine.md](docs/statemachine.md) | `internal/usecase/sample` |
| Both | Both docs above | `users` + `sample` |

## Required workflow

1. Read the relevant doc(s) above.
2. Follow existing naming and layer boundaries (`entity` → `repository` → `usecase` → `transport` → `bootstrap`).
3. Wire new services in `internal/bootstrap/wire.go` and register routes in `internal/transport/apis/router.go`.
4. Add GORM models to `internal/infrastructure/database/postgres/connection.go` `Migrate()` when adding tables.
5. Run `go build ./...` (and `go vet ./...` when possible) before finishing.

## Conventions worth noting

- Domain **entity** packages may be plural (`users`) while **repository** folders are often singular (`user`) — match existing domains.
- Usecase depends on **interfaces**, not postgres/kafka implementations.
- Inject `begin.BeginRepository` into usecases that use DB transactions (`NewUserService`, sample `TransactionManager`).
- Sample is the **end-to-end** state machine example: `POST /samples`, `PUT /samples/:id`.

## Do not

- Put GORM tags on `internal/entity/*`.
- Put status transition logic in HTTP handlers.
- Commit secrets or `.env` files.
