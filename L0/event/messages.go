package event

import (
	"L0/config"
)

type Message interface {
	Key() string
}

type OrderCreatedMessage struct {
	ID   string
	Body Order
}

func (m *OrderCreatedMessage) Key() string {
	return config.MsgCreated
}
