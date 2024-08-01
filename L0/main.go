package main

import (
	"L0/event"
	"L0/schema"
	"fmt"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// общие мысли: организовать параллельный вызов обрабатываемых сервером функций и в принципе
// распараллелить все необходимые процессы

// client (?)
// "простейший интерфейс" - какая-нибудь html-ка XD
// подумать как ограничить объекты , входящие в канал (хитрая организация интерфйесов)

func main() {
	dsn := "host=localhost user=max_db password=max_db dbname=golang_base port=5433"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}

	// Миграция схемы
	err = db.AutoMigrate(&schema.Order{}, &schema.Delivery{}, &schema.Payment{}, &schema.Item{})
	if err != nil {
		log.Fatal("Ошибка при миграции схемы:", err)
	}

	// Пример создания нового заказа
	delivery := schema.Delivery{
		Name:    "John Doe",
		Phone:   "+1234567890",
		Zip:     "123456",
		City:    "Somewhere",
		Address: "123 Main St",
		Region:  "Some Region",
		Email:   "john@example.com",
	}

	payment := schema.Payment{
		Transaction:  "trans123",
		RequestID:    "req123",
		Currency:     "USD",
		Provider:     "provider1",
		Amount:       100.5,
		PaymentDt:    1627849820,
		Bank:         "Some Bank",
		DeliveryCost: 10.0,
		GoodsTotal:   90.5,
		CustomFee:    0.0,
	}

	order := schema.Order{
		OrderUID:        "order123",
		TrackNumber:     "track123",
		Entry:           "entry1",
		Delivery:        delivery,
		Payment:         payment,
		Locale:          "en",
		InternalSig:     "sig1",
		CustomerID:      "customer123",
		DeliveryService: "service1",
		ShardKey:        "shard1",
		SMID:            1,
		DateCreated:     time.Now(),
		OofShard:        "oof1",
		Items: []schema.Item{
			{ChrtID: 1, TrackNumber: "track1", Price: 50.0, RID: "rid1", Name: "item1", Sale: 0, Size: "M", TotalPrice: 50.0, NmID: 1, Brand: "brand1", Status: 1},
			{ChrtID: 2, TrackNumber: "track2", Price: 40.5, RID: "rid2", Name: "item2", Sale: 0, Size: "L", TotalPrice: 40.5, NmID: 2, Brand: "brand2", Status: 1},
		},
	}

	// Сохранение заказа в базу данных
	result := db.Create(&order)
	if result.Error != nil {
		log.Fatal("Ошибка при создании заказа:", result.Error)
	}

	// Вывод данных о созданном заказе
	fmt.Printf("Заказ успешно создан: %+v\n", order)

	es, err := event.NewNats(nats.DefaultURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer es.Close()

	store := &event.InMemoryStore{Orders: make(map[string]schema.Order)}

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
	//http.HandleFunc("/orders", event.GetAllOrdersHandler(store))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
