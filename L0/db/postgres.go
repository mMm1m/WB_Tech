package db

import (
	"L0/schema"
	_ "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgres() (*PostgresRepository, error) {
	dsn := "host=localhost user=max_db password=max_db dbname=golang_base port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}

	err = db.AutoMigrate(&schema.Order{}, &schema.Delivery{}, &schema.Payment{}, &schema.Item{})
	if err != nil {
		log.Fatal("Ошибка при миграции схемы:", err)
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	sqlDB, _ := r.db.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}
}

func (r *PostgresRepository) InsertOrder(order schema.Order) error {
	result := r.db.Create(&order)
	return result.Error
}

func (r *PostgresRepository) GetAllOrders() ([]schema.Order, error) {
	var orders []schema.Order
	result := r.db.Find(&orders)
	return orders, result.Error
}

func (r *PostgresRepository) AlreadyExists(order schema.Order) (bool, error) {
	result := r.db.Where("order_uid = ?", order.OrderUID).First(&order)
	if result.Error == nil {
		fmt.Println("Такой заказ уже существует:", order)
		return false, nil
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Fatal("Ошибка при проверке существования заказа:", result.Error)
		return false, result.Error
	}

	return true, result.Error
}
