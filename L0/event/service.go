package event

import (
	"L0/config"
	"L0/db"
	"L0/schema"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

type InMemoryStore struct {
	mu     sync.Mutex
	Orders map[string]schema.Order
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Orders: make(map[string]schema.Order),
	}
}

func (store *InMemoryStore) GetOrder(es *NatsEventStore) error {
	sub, err := es.Nc.Subscribe(config.PostCluster, func(msg *nats.Msg) {
		restoreDBCache(store)
		orderID := string(msg.Data)
		store.mu.Lock()
		order, exists := store.Orders[orderID]
		store.mu.Unlock()

		var response []byte
		if exists {
			response, _ = json.Marshal(order)
		} else {
			response = []byte("")
		}

		err := es.Nc.Publish(msg.Reply, response)
		if err != nil {
			log.Println("Failed to publish reply:", err)
		}
	})
	if err != nil {
		return err
	}

	es.OrderCreatedSubscription = sub
	return nil
}

func (store *InMemoryStore) AddOrder(es *NatsEventStore) error {
	sub, err := es.Nc.Subscribe(config.GetCluster, func(msg *nats.Msg) {
		restoreDBCache(store)
		var order schema.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		ok, err_ := db.AlreadyExists(order)
		if err_ != nil && ok == true {
			db.InsertOrder(order)
		}
		store.mu.Lock()
		store.Orders[order.OrderUID] = order
		store.mu.Unlock()
	})
	if err != nil {
		return err
	}

	es.OrderCreatedSubscription = sub
	return nil
}

func restoreDBCache(store *InMemoryStore) {
	allOrders, _ := db.GetAllOrders()
	if len(allOrders) != len(store.Orders) {
		for _, order := range allOrders {
			store.Orders[order.OrderUID] = order
		}
	}
}
