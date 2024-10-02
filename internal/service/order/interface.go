package service

import "saga-order/internal/model"

type SagaOrderService interface {
	CreateOrder(createInput model.Order) (model.Order, error)
	RollbackOrder(orderID uint64) error
}
