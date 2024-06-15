package main

import (
	"net/http"
)

var (
	store = NewInMemoryStore()
)

func main() {
	// HTTP handler для GET запроса
	http.HandleFunc("/order", GetOrderHandler)
	// HTTP handler для POST запроса
	http.HandleFunc("/add", AddOrderHandler)
	// Запуск HTTP сервера на порту 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// handle error
	}
}
