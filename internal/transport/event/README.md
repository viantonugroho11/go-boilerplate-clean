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

3. **Runner & bootstrap** — di `router.go`: (1) tambah konstanta nama di `ConsumerNames` dan const (e.g. `ConsumerNameXxx`), (2) tambah fungsi `RunXxx(ctx, cfg, ...deps) (kafka.Consumer, error)`, (3) di `bootstrap/consumer.go` tambah `case event.ConsumerNameXxx:` yang panggil runner tersebut. Flag **single** `-consumer=<nama>`:
   ```bash
   ./consumer -consumer=user
   ./consumer -consumer=order
   ```
   Hanya satu consumer per proses; dependency (e.g. DB) hanya di-init di case yang butuh.

## Contoh yang ada

- **User**: `events.UserCreatedEvent` + `kafka.UserCreatedHandler` (panggil `UserService.GetByID`).
- **Order**: `events.OrderCreatedEvent` + `kafka.OrderCreatedHandler` (contoh, log saja).
