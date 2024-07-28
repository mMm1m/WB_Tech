package event

type Message interface {
	Key() string
}

type OrderCreatedMessage struct {
	ID   string
	Body Order
}

func (m *OrderCreatedMessage) Key() string {
	return "message created"
}
