package main

import (
	"L0/event"
	_ "fmt"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

// общие мысли: организовать параллельный вызов обрабатываемых сервером функций и в принципе
// распараллелить все необходимые процессы

// client (?)
// "простейший интерфейс" - какая-нибудь html-ка XD
// подумать как ограничить объекты , входящие в канал (хитрая организация интерфйесов)

func main() {
	log.Println("Starting...")

	es, err := event.NewNats(nats.DefaultURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer es.Close()

	store := &event.InMemoryStore{Orders: make(map[string]event.Order)}

	err = store.AddOrder(es)
	if err != nil {
		log.Println("Failed to subscribe to order creation events:", err)
		return
	}

	err = store.GetOrder(es)
	if err != nil {
		log.Println("Failed to subscribe to get order requests:", err)
		return
	}
	http.HandleFunc("/order", event.GetOrderHandler(es))
	http.HandleFunc("/order/add", event.AddOrderHandler(es))
	http.HandleFunc("/orders", event.GetAllOrdersHandler(store))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
