package event

import (
	"L0/config"
	"L0/schema"
)

type Message interface {
	Key() string
}

type OrderCreatedMessage struct {
	ID   string
	Body schema.Order
}

func (m *OrderCreatedMessage) Key() string {
	return config.MsgCreated
}
