package event

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetOrderHandler(es *NatsEventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("order_uid")
		if orderID == "" {
			http.Error(w, "Missing order ID", http.StatusBadRequest)
			return
		}

		msg, err := es.nc.Request("order.got", []byte(orderID), 5*time.Second)
		if err != nil {
			http.Error(w, "Failed to get order", http.StatusInternalServerError)
			return
		}

		var order Order
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			http.Error(w, "Failed to decode order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err_ := json.NewEncoder(w).Encode(order)
		if err_ != nil {
			return
		}
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
		}

		w.WriteHeader(http.StatusCreated)
		err_ := json.NewEncoder(w).Encode(order)
		if err_ != nil {
			return
		}
	}
}

/*func GetAllOrdersHandler(store *InMemoryStore) http.HandlerFunc {
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
}*/
