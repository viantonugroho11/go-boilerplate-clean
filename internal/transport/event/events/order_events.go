package events

// OrderCreatedEvent payload untuk event order created (Kafka JSON).
// Contoh event untuk consumer lain; sesuaikan field-nya dengan kebutuhan.
type OrderCreatedEvent struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}
