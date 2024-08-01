package event

import (
	"L0/config"
	"L0/schema"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetOrderHandler(es *NatsEventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}
		//fmt.Print(idealJSONFields)
		orderID := r.URL.Query().Get(config.OrderID)
		if orderID == "" {
			http.Error(w, "Missing order ID", http.StatusBadRequest)
			return
		}

		msg, err := es.Nc.Request(config.PostCluster, []byte(orderID), 5*time.Second)
		if err != nil {
			http.Error(w, "Failed to get order", http.StatusInternalServerError)
			return
		}

		if len(msg.Data) == 0 {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		var order schema.Order
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
		var order schema.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		msg, _ := json.Marshal(order)
		err := es.Nc.Publish(config.GetCluster, msg)
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

func RestoreDBCache() {}

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
