package db

import "L0/schema"

type Repository interface {
	Close()
	InsertOrder(order schema.Order) error
	GetAllOrders() ([]schema.Order, error)
	AlreadyExists(order schema.Order) (bool, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertOrder(order schema.Order) error {
	return impl.InsertOrder(order)
}

func GetAllOrders() ([]schema.Order, error) {
	return impl.GetAllOrders()
}

func AlreadyExists(order schema.Order) (bool, error) {
	return impl.AlreadyExists(order)
}
