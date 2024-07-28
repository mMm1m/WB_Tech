package event

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsEventStore struct {
	Nc                       *nats.Conn
	OrderCreatedSubscription *nats.Subscription
	orderCreatedChan         chan OrderCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{Nc: nc}, nil
}

func (es *NatsEventStore) Close() {
	if es.Nc != nil {
		es.Nc.Close()
	}
	if es.OrderCreatedSubscription != nil {
		if err := es.OrderCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	if es.orderCreatedChan != nil {
		close(es.orderCreatedChan)
	}
}
