package event

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetOrderHandler(store *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		orderUID := r.URL.Query().Get("order_uid")
		if orderUID == "" {
			http.Error(w, "Missing order_uid parameter", http.StatusBadRequest)
			return
		}

		order, ok := store.GetOrder(orderUID)
		if !ok {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}

func AddOrderHandler(es *NatsEventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		msg, _ := json.Marshal(order)
		err := es.nc.Publish("order.created", msg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to publish order", http.StatusInternalServerError)
			return
		} else {
			log.Printf("Publisher  =>  Message: %s\n", order.OrderUID)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func GetAllOrdersHandler(store *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		orders := store.getAllOrders()
		jsonData, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
