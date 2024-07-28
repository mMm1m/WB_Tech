package event

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsEventStore struct {
	nc                       *nats.Conn
	orderCreatedSubscription *nats.Subscription
	orderCreatedChan         chan OrderCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

// Close закрывает соединение и отписывается от событий
func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.orderCreatedSubscription != nil {
		if err := es.orderCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	if es.orderCreatedChan != nil {
		close(es.orderCreatedChan)
	}
}
