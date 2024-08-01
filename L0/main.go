package main

import (
	"L0/db"
	"L0/event"
	"L0/schema"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

func main() {
	repo, err := db.NewPostgres()
	if repo == nil {
		log.Fatal("Database is not defined")
	}
	db.SetRepository(repo)
	defer db.Close()

	es, err := event.NewNats(nats.DefaultURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer es.Close()

	store := &event.InMemoryStore{Orders: make(map[string]schema.Order)}

	err = store.AddOrder(es)
	if err != nil {
		log.Println("Failed to subscribe to add orders requests:", err)
		return
	}

	err = store.GetOrder(es)
	if err != nil {
		log.Println("Failed to subscribe to get order requests:", err)
		return
	}
	http.HandleFunc("/order", event.GetOrderHandler(es))
	http.HandleFunc("/order/add", event.AddOrderHandler(es))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
