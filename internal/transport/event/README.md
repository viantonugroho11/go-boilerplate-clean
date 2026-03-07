# Event consumers (Kafka)

Consumer memakai **go-lib/kafka**; handler di sini, factory di `internal/infrastructure/broker/kafka`.

## Menambah consumer baru

1. **Event payload** — di `events/`  
   Buat struct dengan tag `json` (decode otomatis oleh go-lib).  
   Contoh: `events/order_events.go` → `OrderCreatedEvent`.

2. **Handler** — di `kafka/`  
   Buat struct yang implement `kafka.EventHandler[EventType]`:  
   - `Name() string`  
   - `Handle(ctx, evt, headers...) kafka.Progress`  
   Inject usecase di constructor, panggil usecase di `Handle`.  
   Contoh: `kafka/order_consumer_handler.go` → `OrderCreatedHandler`.

3. **Runner & main** — di `kafka/runner.go` tambah `RunXxx(ctx, cfg, ...deps) (kafka.Consumer, error)`; di `cmd/consumer/main.go` tambah `case "nama":` yang panggil runner tersebut. Flag **single** `-consumer=<nama>`:
   ```bash
   ./consumer -consumer=user
   ./consumer -consumer=order
   ```
   Hanya satu consumer per proses; dependency (e.g. DB) hanya di-init di case yang butuh.

## Contoh yang ada

- **User**: `events.UserCreatedEvent` + `kafka.UserCreatedHandler` (panggil `UserService.GetByID`).
- **Order**: `events.OrderCreatedEvent` + `kafka.OrderCreatedHandler` (contoh, log saja).
