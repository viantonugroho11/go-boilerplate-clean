# Go Boilerplate – Clean Architecture

A structured Go boilerplate (clean-architecture inspired) using:
- Echo (HTTP server)
- GORM + PostgreSQL (ORM & database)
- Kafka (producer & consumer) via `github.com/IBM/sarama`
- Redis client via `github.com/redis/go-redis/v9`
- Viper for configuration (file + environment override)

Goals: clear folder layout, separation of concerns (usecase, repository, transport, infrastructure), and explicit dependency wiring.

## Prerequisites
- Go (per `go.mod`, Go 1.23+)
- PostgreSQL
- Kafka broker (+ Zookeeper for the provided compose)
- Redis

## Folder Structure
```
cmd/
  app/                 # application entrypoint (main)
internal/
  config/              # configuration loader (Viper) & DSN helpers
  entity/              # domain entities (e.g., User)
  infrastructure/
    database/
      postgres/        # GORM connection & migration
        connection.go
        migrate/       # optional SQL examples
    broker/
      kafka/           # Kafka producer & consumer wrappers
    cache/
      redis/           # Redis client initialization
  repository/
    user/
      model/           # GORM model mappings
      postgres/        # repository implementation with GORM
      user_repository.go  # repository interface
  transport/
    apis/              # HTTP router
    http/
      dto/             # request/response DTOs
      handler/         # Echo handlers
    event/
      kafka/           # consumer runner & example handlers
  usecase/             # application/service layer
configs/
  config.yaml          # default configuration (overridable by env)
Dockerfile
docker-compose.yml
```

## Configuration (Viper)
Configuration is read from `configs/config.yaml` and can be overridden by environment variables (Viper’s `AutomaticEnv`).

HTTP:
- `PORT` (default: `8080`)

PostgreSQL (choose one approach):
- `DATABASE_URL` (e.g., `postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable`)
- or separate fields:
  - `DB_HOST` (default: `127.0.0.1`)
  - `DB_PORT` (default: `5432`)
  - `DB_USER` (default: `postgres`)
  - `DB_PASSWORD` (default: empty)
  - `DB_NAME` (default: `appdb`)
  - `DB_SSLMODE` (default: `disable`)

Kafka:
- `KAFKA_BROKERS` (default: `127.0.0.1:9092`, comma-separated)
- `KAFKA_CLIENT_ID` (default: `go-boilerplate-clean`)
- `KAFKA_GROUP_ID` (default: `go-boilerplate-clean-group`)
- `KAFKA_TOPIC` (default: `user-events`)

Redis:
- `REDIS_ADDR` (default: `127.0.0.1:6379`)
- `REDIS_PASSWORD` (default: empty)
- `REDIS_DB` (default: `0`)

See `configs/config.yaml` for a docker-friendly baseline.

## Run Locally
1) Install dependencies (optional; `go mod tidy` will fetch them):
```bash
go get gorm.io/gorm gorm.io/driver/postgres
go get github.com/labstack/echo/v4
go get github.com/IBM/sarama
go get github.com/redis/go-redis/v9
go get github.com/google/uuid
go get github.com/spf13/viper
go mod tidy
```

2) Set environment variables (example):
```bash
export PORT=8080
export DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"
export KAFKA_BROKERS="127.0.0.1:9092"
export KAFKA_CLIENT_ID="go-boilerplate-clean"
export KAFKA_GROUP_ID="go-boilerplate-clean-group"
export KAFKA_TOPIC="user-events"
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_DB="0"
```

3) Run the server:
```bash
go run ./cmd/app
```

At startup the app will:
- Initialize a GORM connection to PostgreSQL
- AutoMigrate the `users` table
- Start Echo HTTP server
- Initialize Redis client
- Initialize Kafka producer & consumer (consumer runs in background)

## Docker & Compose
Build the image:
```bash
docker build -t go-boilerplate-clean:latest .
```

Run with compose (app + Postgres + Zookeeper + Kafka + Redis):
```bash
docker compose up -d --build
```

Exposed ports:
- App: `8080`
- Postgres: `5432`
- Kafka: `9092`
- Redis: `6379`

## HTTP Endpoints
Base URL: `http://localhost:${PORT}`

Healthcheck:
```bash
GET /healthz
```

User CRUD:
- `POST /users`
  - Body:
    ```json
    { "name": "Jane", "email": "jane@example.com" }
    ```
- `GET /users`
- `GET /users/:id`
- `PUT /users/:id`
  - Body:
    ```json
    { "name": "Jane Updated", "email": "jane.updated@example.com" }
    ```
- `DELETE /users/:id`

Example cURL:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane","email":"jane@example.com"}'
```

## Kafka
Producer:
- Wrapper at `internal/infrastructure/broker/kafka/producer.go`, use `Publish(ctx, topic, key, value)`.

Consumer:
- Wrapper at `internal/infrastructure/broker/kafka/consumer.go` using a consumer group.
- Registration & start in `internal/transport/event/kafka/consumer_runner.go` (example handler logs messages).
- Wired in `cmd/app/main.go` via group and topic.

Notes:
- Ensure `KAFKA_BROKERS` points to a running broker.
- Replace `ExampleHandler` to call real usecases as needed.

## Redis
Client:
- Initialization in `internal/infrastructure/cache/redis/client.go`.
- Wired in `cmd/app/main.go`. You can inject it into layers for caching, rate limiting, etc.

## Database
ORM:
- GORM model for `users` in `internal/repository/user/model/user.go`.
- AutoMigrate runs on startup.

Optional SQL:
- Example file at `internal/infrastructure/database/postgres/migrate/init_users.sql`.

## Repository & Usecase
- Repository interface: `internal/repository/user/user_repository.go`
- Postgres (GORM) implementation: `internal/repository/user/postgres/repository.go`
- Usecase: `internal/usecase/user_usecase.go`
- HTTP handlers (Echo): `internal/transport/http/handler/user_handler.go`
- HTTP DTOs: `internal/transport/http/dto/`

## Architecture Notes
- `usecase` depends only on the repository interface
- `repository/*/postgres` adapts the interface with GORM
- `transport/http` (Echo) calls usecases
- `infrastructure` holds I/O details (database, kafka, redis)

## Troubleshooting
- Build dependency issues: run `go mod tidy`
- Postgres connection issues: verify `DATABASE_URL` or `DB_*` variables
- Kafka broker unavailable: verify `KAFKA_BROKERS` and broker status
- Redis connection refused: verify `REDIS_ADDR` and Redis status

## License
MIT. Feel free to use and modify.

# go-boilerplate-clean