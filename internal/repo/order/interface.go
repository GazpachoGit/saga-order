package repo

import "saga-order/internal/model"

type OrderRepo interface {
	CreateOrder(createInput model.Order) (model.Order, error)
	DeleteOrder(orderID uint64) error
}
