package event

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type Order struct {
	OrderUID        string    `json:"order_uid"`
	TrackNumber     string    `json:"track_number"`
	Entry           string    `json:"entry"`
	Delivery        Delivery  `json:"delivery"`
	Payment         Payment   `json:"payment"`
	Items           []Item    `json:"items"`
	Locale          string    `json:"locale"`
	InternalSig     string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id"`
	DeliveryService string    `json:"delivery_service"`
	ShardKey        string    `json:"shardkey"`
	SMID            int       `json:"sm_id"`
	DateCreated     time.Time `json:"date_created"`
	OofShard        string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	RID         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

type InMemoryStore struct {
	mu     sync.Mutex
	Orders map[string]Order
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Orders: make(map[string]Order),
	}
}

func (store *InMemoryStore) GetOrder(es *NatsEventStore) error {
	sub, err := es.nc.Subscribe("order.got", func(msg *nats.Msg) {
		orderID := string(msg.Data)
		store.mu.Lock()
		order, exists := store.Orders[orderID]
		store.mu.Unlock()

		var response []byte
		if exists {
			response, _ = json.Marshal(order)
		} else {
			response = []byte("{}")
		}

		err := es.nc.Publish(msg.Reply, response)
		if err != nil {
			log.Println("Failed to publish reply:", err)
		}
	})
	if err != nil {
		return err
	}

	es.orderCreatedSubscription = sub
	return nil
}

func (store *InMemoryStore) AddOrder(es *NatsEventStore) error {
	sub, err := es.nc.Subscribe("order.created", func(msg *nats.Msg) {
		var order Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		store.mu.Lock()
		store.Orders[order.OrderUID] = order
		store.mu.Unlock()
		print(len(store.Orders))
	})
	if err != nil {
		return err
	}

	es.orderCreatedSubscription = sub
	return nil
}

/*func (store *InMemoryStore) getAllOrders() []Order {
	var orders []Order
	for _, order := range store.Orders {
		orders = append(orders, order)
	}
	return orders
}*/
