package event

import (
	"L0/config"
	"L0/schema"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"
)

func GetOrderHandler(es *NatsEventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

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

func validateStruct(data []byte, model interface{}) error {
	var unnecessaryFields = make(map[string]bool)
	unnecessaryFields[""] = true
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return err
	}

	modelType := reflect.TypeOf(model).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		exists := unnecessaryFields[field.Tag.Get("json")]
		if exists == false {
			if _, ok := jsonObj[field.Tag.Get("json")]; !ok {
				return fmt.Errorf("missing field: %s", field.Name)
			}
		}
	}
	return nil
}

func AddOrderHandler(es *NatsEventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var order schema.Order
		data, _ := io.ReadAll(r.Body)
		if err := validateStruct(data, &order); err != nil {
			http.Error(w, "Incorrect transform", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&order); err != nil {
			fmt.Print(err)
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
