package main

import (
	"encoding/json"
	"net/http"
)

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	// несоответствие методу
	if r.Method != http.MethodGet {
		// временная обработка ошибок
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// проверка работоспособности сервиса
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// проверка сервиса на валидность
	responseCode, _ := checkServiceAvailability(r.URL.String(), client)
	if responseCode == http.StatusServiceUnavailable {
		// handle that service in unavailable
		http.Error(w, "Unavailable service", http.StatusServiceUnavailable)
		return
	}

	// без данного идентификатора поиск бесполезен
	orderUID := r.URL.Query().Get("order_uid")
	if orderUID == "" {
		// временная обработка ошибок
		http.Error(w, "Missing order_uid parameter", http.StatusBadRequest)
		return
	}

	// достаём заказ по индексу из in-memory / базы данных
	order, ok := store.GetOrder(orderUID)
	if !ok {
		// временная обработка ошибок
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// преобразуем в json-формат и обрабатываем ошибку при неудачном преобразовании
	jsonData, err := json.Marshal(order)
	if err != nil {
		// временная обработка ошибок
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// настройка NATS-STREAMING и закидываем текущий json туда
	var order Order
	// получаем json из тела запроса и преобразуем в Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// реализовать временного клиента с помощью контекста,
	// чтобы переадресация заканчивалась спустя n секунд (?)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// проверка сервиса на валидность
	responseCode, _ := checkServiceAvailability(r.URL.String(), client)
	if responseCode == http.StatusServiceUnavailable {
		// handle that service in unavailable
		return
	}

	// читаем из nats и закидываем все объекты в базу данных,
	// и если записей в нем больше одной - выполняем RestoreDBCache

	// иначе просто сохраняем одну запись в in-memory storage (п)
	store.AddOrder(order)

	w.WriteHeader(http.StatusOK)
}

func checkServiceAvailability(url string, client *http.Client) (int, bool) {
	resp, err := client.Get(url)
	if err != nil {
		// handle error
		return http.StatusServiceUnavailable, false
	}
	defer resp.Body.Close()
	// Проверка кода состояния
	if resp.StatusCode/100 == 3 {
		checkServiceAvailability(resp.Header.Get("Location"), client)
	}
	if resp.StatusCode/100 == 2 {
		return resp.StatusCode, true
	}
	return http.StatusServiceUnavailable, false
}

// восстановление данных в in-memory из базы данных (вопрос с сигнатурой)
func RestoreDBCache() {}
